package vault

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azsecrets"
)

// AzureVaultClient implements the VaultClient interface for a Azure Key Vault
type AzureVaultClient struct {
	client        *azsecrets.Client
	azureVaultURL string
}

// validateEnvironmentVariables validates the required environment variables
func (v *AzureVaultClient) validateEnvironmentVariables() error {
	vaultURL := os.Getenv("AZURE_KEYVAULT_URL")
	if vaultURL == "" {
		return fmt.Errorf("AZURE_KEYVAULT_URL environment variable is required")
	}

	v.azureVaultURL = vaultURL

	return nil
}

// InitializeClient initializes the Azure client
func (v *AzureVaultClient) InitializeClient() error {
	err := v.validateEnvironmentVariables()
	if err != nil {
		return err
	}

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return fmt.Errorf("failed to acquire Azure credentials: %v", err)
	}
	// Establish a connection to the Key Vault client
	client, err := azsecrets.NewClient(v.azureVaultURL, cred, nil)

	v.client = client

	return nil
}

// GetSecrets retrieves secrets from the Azure Key Vault instance
func (v *AzureVaultClient) GetSecrets() ([]VaultSecret, error) {
	var secrets []VaultSecret

	pager := v.client.NewListSecretsPager(nil)
	for pager.More() {
		page, err := pager.NextPage(context.Background())
		if err != nil {
			return nil, err
		}
		for _, secret := range page.Value {
			pullSecret, err := v.client.GetSecret(context.Background(), secret.ID.Name(), secret.ID.Version(), &azsecrets.GetSecretOptions{})
			if err != nil {
				return nil, err
			}
			repositories := ""
			if val, ok := pullSecret.Tags["repositories"]; ok && val != nil {
				repositories = *val
			}

			newSecret := VaultSecret{
				Name:         pullSecret.ID.Name(),
				Value:        *pullSecret.Value,
				Type:         *pullSecret.Tags["type"],
				Visibility:   *pullSecret.Tags["visibility"],
				Repositories: strings.Split(repositories, ","),
			}
			secrets = append(secrets, newSecret)
		}
	}

	return secrets, nil
}
