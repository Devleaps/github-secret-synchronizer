package github

import (
	"context"
	"encoding/base64"
	"errors"
	"net/http"
	"os"
	"strconv"

	"github.com/Devleaps/github-secret-synchronizer/internal/vault"
	"github.com/bradleyfalzon/ghinstallation/v2"
	"github.com/google/go-github/v64/github"
	"golang.org/x/crypto/nacl/box"
)

type GitHubWrapper struct {
	client *github.Client

	appID          int64
	installationID int64
	privateKey     []byte
	orgName        string

	publicKey string
	keyID     string
}

// NewClient creates a new GitHub client
func (g *GitHubWrapper) NewClient() error {
	appID, installationID, privateKey, orgName, err := getGitHubAppCredentials()
	if err != nil {
		return err
	}
	g.appID = appID
	g.installationID = installationID
	g.privateKey = privateKey
	g.orgName = orgName

	itr, err := ghinstallation.New(http.DefaultTransport, g.appID, g.installationID, g.privateKey)
	if err != nil {
		return err
	}

	client := github.NewClient(&http.Client{Transport: itr})
	g.client = client

	publicKey, keyID, err := g.getOrgPublicKey()
	if err != nil {
		return err
	}
	g.publicKey = publicKey
	g.keyID = keyID
	return nil
}

// PushEncryptedSecret pushes the encrypted secret to GitHub
func (g *GitHubWrapper) PushSecret(secret vault.VaultSecret) error {
	encryptedSecret, err := encryptSecret(secret.Value, g.publicKey)
	if err != nil {
		return err
	}
	var ghSecret *github.EncryptedSecret
	if secret.Visibility == vault.SELECTED_VISIBILITY {
		repositoryIds, err := g.retrieveRepositoriesIds(secret.Repositories)
		if err != nil {
			return err
		}
		ghSecret = &github.EncryptedSecret{
			Name:                  secret.Name,
			KeyID:                 g.keyID,
			EncryptedValue:        encryptedSecret,
			Visibility:            secret.Visibility,
			SelectedRepositoryIDs: repositoryIds,
		}
	} else {
		ghSecret = &github.EncryptedSecret{
			Name:           secret.Name,
			KeyID:          g.keyID,
			EncryptedValue: encryptedSecret,
			Visibility:     secret.Visibility,
		}
	}

	ctx := context.Background()
	_, err = g.client.Actions.CreateOrUpdateOrgSecret(ctx, g.orgName, ghSecret)

	return err
}

// PushVariable pushes the variable to GitHub
func (g *GitHubWrapper) PushVariable(variable vault.VaultSecret) error {
	var ghVariable *github.ActionsVariable
	if variable.Visibility == vault.SELECTED_VISIBILITY {
		repositoryIds, err := g.retrieveRepositoriesIds(variable.Repositories)
		if err != nil {
			return err
		}
		ghVariable = &github.ActionsVariable{
			Name:                  variable.Name,
			Value:                 variable.Value,
			Visibility:            &variable.Visibility,
			SelectedRepositoryIDs: &repositoryIds,
		}
	} else {
		ghVariable = &github.ActionsVariable{
			Name:       variable.Name,
			Value:      variable.Value,
			Visibility: &variable.Visibility,
		}
	}
	err := g.createOrUpdateVariable(ghVariable)

	return err
}

func (g *GitHubWrapper) retrieveRepositoriesIds(repositories []string) (github.SelectedRepoIDs, error) {
	var repositoryIds []int64
	for _, repository := range repositories {
		repo, _, err := g.client.Repositories.Get(context.Background(), g.orgName, repository)
		if err != nil {
			return nil, err
		}
		repositoryIds = append(repositoryIds, repo.GetID())
	}
	return repositoryIds, nil
}

// createOrUpdateVariable creates or updates a variable in the organization
func (g *GitHubWrapper) createOrUpdateVariable(variable *github.ActionsVariable) error {
	ctx := context.Background()
	_, res, err := g.client.Actions.GetOrgVariable(ctx, g.orgName, variable.Name)
	if res.StatusCode == 404 {
		_, err = g.client.Actions.CreateOrgVariable(ctx, g.orgName, variable)
	} else if res.StatusCode == 200 {
		_, err = g.client.Actions.UpdateOrgVariable(ctx, g.orgName, variable)
	} else {
		return errors.New("Unable to get variable, status code: " + strconv.Itoa(res.StatusCode))
	}

	return err
}

// getOrgPublicKey retrieves the organization public key
func (g *GitHubWrapper) getOrgPublicKey() (string, string, error) {
	ctx := context.Background()
	publicKey, _, err := g.client.Actions.GetOrgPublicKey(ctx, g.orgName)
	if err != nil {
		return "", "", err
	}
	return publicKey.GetKey(), publicKey.GetKeyID(), nil
}

// ecryptSecret encrypts a secret value using the public key
func encryptSecret(secretValue, publicKey string) (string, error) {
	publicKeyBytes, err := base64.StdEncoding.DecodeString(publicKey)
	if err != nil {
		return "", err
	}

	var publicKeyArray [32]byte
	copy(publicKeyArray[:], publicKeyBytes)

	encrypted, err := box.SealAnonymous(nil, []byte(secretValue), &publicKeyArray, nil)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(encrypted), nil
}

// getGitHubAppCredentials retrieves the GitHub App credentials from the environment
func getGitHubAppCredentials() (int64, int64, []byte, string, error) {
	appID, exists := os.LookupEnv("GITHUB_APP_ID")
	if !exists {
		return 0, 0, nil, "", errors.New("GITHUB_APP_ID environment variable is required")
	}

	installationID, exists := os.LookupEnv("GITHUB_INSTALLATION_ID")
	if !exists {
		return 0, 0, nil, "", errors.New("GITHUB_INSTALLATION_ID environment variable is required")
	}

	privateKey, exists := os.LookupEnv("GITHUB_PRIVATE_KEY")
	if !exists {
		return 0, 0, nil, "", errors.New("GITHUB_PRIVATE_KEY environment variable is required")
	}

	orgName, exists := os.LookupEnv("GITHUB_ORG_NAME")
	if !exists {
		return 0, 0, nil, "", errors.New("GITHUB_ORG_NAME environment variable is required")
	}

	parsedAppID, err := strconv.ParseInt(appID, 10, 64)
	if err != nil {
		return 0, 0, nil, "", errors.New("Failed to parse GITHUB_APP_ID")
	}

	parsedInstallationID, err := strconv.ParseInt(installationID, 10, 64)
	if err != nil {
		return 0, 0, nil, "", errors.New("Failed to parse GITHUB_INSTALLATION_ID")
	}

	return parsedAppID, parsedInstallationID, []byte(privateKey), orgName, nil
}
