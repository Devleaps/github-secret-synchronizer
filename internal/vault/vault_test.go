package vault

import (
	"os"
	"testing"
)

func TestHandleDefaults(t *testing.T) {
	t.Run("sets default visibility, type, and repositories", func(t *testing.T) {
		os.Setenv("DEFAULT_VISIBILITY", "private")
		os.Setenv("DEFAULT_TYPE", "variable")
		os.Setenv("DEFAULT_REPOSITORIES", "repo1,repo2")

		secrets := []VaultSecret{
			{Name: "secret1", Value: "value1"},
			{Name: "secret2", Value: "value2", Visibility: "all"},
			{Name: "secret3", Value: "value3", Type: "secret"},
		}

		err := HandleDefaults(secrets)
		if err != nil {
			t.Fatalf("Expected no error, but got %v", err)
		}

		if secrets[0].Visibility != "private" {
			t.Errorf("Expected visibility to be 'private', but got '%s'", secrets[0].Visibility)
		}
		if secrets[0].Type != "variable" {
			t.Errorf("Expected type to be 'variable', but got '%s'", secrets[0].Type)
		}
		if len(secrets[0].Repositories) != 2 || secrets[0].Repositories[0] != "repo1" || secrets[0].Repositories[1] != "repo2" {
			t.Errorf("Expected repositories to be ['repo1', 'repo2'], but got %v", secrets[0].Repositories)
		}

		if secrets[1].Visibility != "all" {
			t.Errorf("Expected visibility to be 'all', but got '%s'", secrets[1].Visibility)
		}
		if secrets[1].Type != "variable" {
			t.Errorf("Expected type to be 'variable', but got '%s'", secrets[1].Type)
		}

		if secrets[2].Visibility != "private" {
			t.Errorf("Expected visibility to be 'private', but got '%s'", secrets[2].Visibility)
		}
		if secrets[2].Type != "secret" {
			t.Errorf("Expected type to be 'secret', but got '%s'", secrets[2].Type)
		}
	})

	t.Run("returns error for invalid default visibility", func(t *testing.T) {
		os.Setenv("DEFAULT_VISIBILITY", "invalid")

		secrets := []VaultSecret{
			{Name: "secret1", Value: "value1"},
		}

		err := HandleDefaults(secrets)
		if err == nil {
			t.Fatalf("Expected error, but got nil")
		}
	})

	t.Run("returns error for invalid default type", func(t *testing.T) {
		os.Setenv("DEFAULT_VISIBILITY", "private")
		os.Setenv("DEFAULT_TYPE", "invalid")

		secrets := []VaultSecret{
			{Name: "secret1", Value: "value1"},
		}

		err := HandleDefaults(secrets)
		if err == nil {
			t.Fatalf("Expected error, but got nil")
		}
	})
}
func TestInitializeVault(t *testing.T) {
	t.Run("returns JSONVaultClient when VAULT_TYPE is json", func(t *testing.T) {
		os.Setenv("VAULT_TYPE", "json")
		defer os.Unsetenv("VAULT_TYPE")

		client, err := InitializeVault()
		if err != nil {
			t.Fatalf("Expected no error, but got %v", err)
		}

		if _, ok := client.(*JSONVaultClient); !ok {
			t.Errorf("Expected client to be of type *JSONVaultClient, but got %T", client)
		}
	})

	t.Run("returns YAMLVaultClient when VAULT_TYPE is yaml", func(t *testing.T) {
		os.Setenv("VAULT_TYPE", "yaml")
		defer os.Unsetenv("VAULT_TYPE")

		client, err := InitializeVault()
		if err != nil {
			t.Fatalf("Expected no error, but got %v", err)
		}

		if _, ok := client.(*YAMLVaultClient); !ok {
			t.Errorf("Expected client to be of type *YAMLVaultClient, but got %T", client)
		}
	})

	t.Run("returns AzureVaultClient when VAULT_TYPE is azure", func(t *testing.T) {
		os.Setenv("VAULT_TYPE", "azure")
		defer os.Unsetenv("VAULT_TYPE")

		client, err := InitializeVault()
		if err != nil {
			t.Fatalf("Expected no error, but got %v", err)
		}

		if _, ok := client.(*AzureVaultClient); !ok {
			t.Errorf("Expected client to be of type *AzureVaultClient, but got %T", client)
		}
	})

	t.Run("returns AWSVaultClient when VAULT_TYPE is aws", func(t *testing.T) {
		os.Setenv("VAULT_TYPE", "aws")
		defer os.Unsetenv("VAULT_TYPE")

		client, err := InitializeVault()
		if err != nil {
			t.Fatalf("Expected no error, but got %v", err)
		}

		if _, ok := client.(*AWSVaultClient); !ok {
			t.Errorf("Expected client to be of type *AWSVaultClient, but got %T", client)
		}
	})

	t.Run("returns error when VAULT_TYPE is invalid", func(t *testing.T) {
		os.Setenv("VAULT_TYPE", "invalid")
		defer os.Unsetenv("VAULT_TYPE")

		_, err := InitializeVault()
		if err == nil {
			t.Fatalf("Expected error, but got nil")
		}

		expectedError := "VAULT_TYPE environment variable is required and must be a valid vault type"
		if err.Error() != expectedError {
			t.Errorf("Expected error message to be '%s', but got '%s'", expectedError, err.Error())
		}
	})
}
