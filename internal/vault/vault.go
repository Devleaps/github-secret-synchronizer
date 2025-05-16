package vault

import (
	"errors"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/rs/zerolog/log"
)

const (
	// All types of visibility
	ALL_VISIBILITY      = "all"
	PRIVATE_VISIBILITY  = "private"
	SELECTED_VISIBILITY = "selected"

	// Typings for secrets/variables in GH
	SECRET_TYPE   = "secret"
	VARIABLE_TYPE = "variable"

	// Defaults for visibility and types
	DEFAULT_TYPE       = SECRET_TYPE
	DEFAULT_VISIBILITY = ALL_VISIBILITY
)

// VaultClient defines the interface for interacting with different vaults
type VaultClient interface {
	InitializeClient() error
	GetSecrets() ([]VaultSecret, error)
}

// VaultSecret represents a secret to be synced to GitHub
type VaultSecret struct {
	Name         string   `json:"name"`
	Value        string   `json:"value"`
	Type         string   `json:"type"`
	Visibility   string   `json:"visibility"`
	Repositories []string `json:"repositories"`
}

// InitializeVault initializes the vault client based on the VAULT_TYPE environment variable
func InitializeVault() (VaultClient, error) {
	vaultType := os.Getenv("VAULT_TYPE")
	switch vaultType {
	case "json":
		return &JSONVaultClient{}, nil
	case "yaml":
		return &YAMLVaultClient{}, nil
	case "azure":
		return &AzureVaultClient{}, nil
	case "aws":
		return &AWSVaultClient{}, nil
	default:
		return nil, errors.New("VAULT_TYPE environment variable is required and must be a valid vault type")
	}
}

// HandleDefaults sets the default values for visibility, type and repositories
func HandleDefaults(secrets []VaultSecret) error {
	//TODO: Make this a one-time operation on boot
	defaultVisibility, err := acquireDefaultVisibility()
	if err != nil {
		return err
	}

	defaultType, err := acquireDefaultType()
	if err != nil {
		return err
	}

	defaultRepositories := acquireDefaultRepositories()

	for i := range secrets {
		secrets[i].Name = formatSecretName(secrets[i].Name)
		if secrets[i].Visibility == "" {
			secrets[i].Visibility = defaultVisibility
		}
		if secrets[i].Type == "" {
			secrets[i].Type = defaultType
		}

		if len(secrets[i].Repositories) == 0 {
			secrets[i].Repositories = defaultRepositories
		}
	}
	return nil
}

// acquireDefaultVisibility returns the default visibility value
func acquireDefaultVisibility() (string, error) {
	defaultVisibility := os.Getenv("DEFAULT_VISIBILITY")
	if defaultVisibility == "" {
		defaultVisibility = DEFAULT_VISIBILITY
	}

	if !slices.Contains([]string{ALL_VISIBILITY, PRIVATE_VISIBILITY, SELECTED_VISIBILITY}, defaultVisibility) {
		return "", fmt.Errorf("invalid visibility value: %s. Please pass one of %s, %s or %s", defaultVisibility, ALL_VISIBILITY, PRIVATE_VISIBILITY, SELECTED_VISIBILITY)
	}
	return defaultVisibility, nil
}

// acquireDefaultType returns the default type value
func acquireDefaultType() (string, error) {
	defaultType := os.Getenv("DEFAULT_TYPE")
	if defaultType == "" {
		defaultType = DEFAULT_TYPE
	}

	if !slices.Contains([]string{SECRET_TYPE, VARIABLE_TYPE}, defaultType) {
		return "", fmt.Errorf("invalid type value: %s. Please pass one of %s or %s", defaultType, SECRET_TYPE, VARIABLE_TYPE)
	}
	return defaultType, nil
}

// acquireDefaultRepositories returns the default repositories value
func acquireDefaultRepositories() []string {
	defaultRepositoriesString := os.Getenv("DEFAULT_REPOSITORIES")

	defaultRepositories := strings.Split(defaultRepositoriesString, ",")

	if len(defaultRepositories) == 0 {
		log.Warn().Msg("DEFAULT_REPOSITORIES is unset")
	}
	return defaultRepositories
}

// formatSecretName formats the secret name to be used as an environment variable
func formatSecretName(secretName string) string {
	return strings.ToUpper(strings.ReplaceAll(strings.ReplaceAll(secretName, " ", "_"), "-", "_"))
}
