package github

import (
	"context"
	"encoding/base64"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/bradleyfalzon/ghinstallation/v2"
	"github.com/google/go-github/v64/github"
	"golang.org/x/crypto/nacl/box"
)

// Client defines the interface for interacting with GitHub
type Client interface {
	NewClient() (*github.Client, error)
	GetOrgPublicKey(client *github.Client) (string, string, error)
	EncryptSecret(secretValue, publicKey string) (string, error)
	PushEncryptedSecret(client *github.Client, secretName, encryptedSecret, keyID string) error
}

// GitHubService implements the Client interface
type GitHubService struct {
	appID          int64
	installationID int64
	privateKey     []byte
	orgName        string
}

// NewGitHubService creates a new GitHubService
func NewGitHubService() *GitHubService {
	appID, installationID, privateKey, orgName := getGitHubAppCredentials()
	return &GitHubService{
		appID:          appID,
		installationID: installationID,
		privateKey:     privateKey,
		orgName:        orgName,
	}
}

// NewClient creates a new GitHub client
func (s *GitHubService) NewClient() (*github.Client, error) {
	itr, err := ghinstallation.New(http.DefaultTransport, s.appID, s.installationID, s.privateKey)
	if err != nil {
		return nil, err
	}

	client := github.NewClient(&http.Client{Transport: itr})
	return client, nil
}

// GetOrgPublicKey retrieves the organization public key
func (s *GitHubService) GetOrgPublicKey(client *github.Client) (string, string, error) {
	ctx := context.Background()
	publicKey, _, err := client.Actions.GetOrgPublicKey(ctx, s.orgName)
	if err != nil {
		return "", "", err
	}
	return publicKey.GetKey(), publicKey.GetKeyID(), nil
}

// EncryptSecret encrypts a secret value using the public key
func (s *GitHubService) EncryptSecret(secretValue, publicKey string) (string, error) {
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

// PushEncryptedSecret pushes the encrypted secret to GitHub
func (s *GitHubService) PushEncryptedSecret(client *github.Client, secretName, encryptedSecret, keyID string, visibility string) error {
	ctx := context.Background()
	secret := &github.EncryptedSecret{
		Name:           secretName,
		KeyID:          keyID,
		EncryptedValue: encryptedSecret,
		Visibility:     visibility,
	}
	_, err := client.Actions.CreateOrUpdateOrgSecret(ctx, s.orgName, secret)
	return err
}

// getGitHubAppCredentials retrieves the GitHub App credentials from the environment
func getGitHubAppCredentials() (int64, int64, []byte, string) {
	appID, exists := os.LookupEnv("GITHUB_APP_ID")
	if !exists {
		log.Fatal("GITHUB_APP_ID environment variable is required")
	}

	installationID, exists := os.LookupEnv("GITHUB_INSTALLATION_ID")
	if !exists {
		log.Fatal("GITHUB_INSTALLATION_ID environment variable is required")
	}

	privateKey, exists := os.LookupEnv("GITHUB_PRIVATE_KEY")
	if !exists {
		log.Fatal("GITHUB_PRIVATE_KEY environment variable is required")
	}

	orgName, exists := os.LookupEnv("GITHUB_ORG_NAME")
	if !exists {
		log.Fatal("GITHUB_ORG_NAME environment variable is required")
	}

	parsedAppID, err := strconv.ParseInt(appID, 10, 64)
	if err != nil {
		log.Fatalf("failed to parse GITHUB_APP_ID: %v", err)
	}

	parsedInstallationID, err := strconv.ParseInt(installationID, 10, 64)
	if err != nil {
		log.Fatalf("failed to parse GITHUB_INSTALLATION_ID: %v", err)
	}

	return parsedAppID, parsedInstallationID, []byte(privateKey), orgName
}
