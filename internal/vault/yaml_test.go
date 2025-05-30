package vault

import (
	"os"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestYAMLInitializeClient(t *testing.T) {
	t.Run("sets filePath from environment variable", func(t *testing.T) {
		expectedPath := "custom_secrets.yaml"
		if err := os.Setenv("YAML_VAULT_FILE_PATH", expectedPath); err != nil {
			t.Fatalf("Failed to set environment variable: %v", err)
		}
		defer func() {
			if err := os.Unsetenv("YAML_VAULT_FILE_PATH"); err != nil {
				t.Fatalf("Failed to unset environment variable: %v", err)
			}
		}()

		client := &YAMLVaultClient{}
		err := client.InitializeClient()
		if err != nil {
			t.Fatalf("Expected no error, but got %v", err)
		}

		if client.filePath != expectedPath {
			t.Errorf("Expected filePath to be '%s', but got '%s'", expectedPath, client.filePath)
		}
	})

	t.Run("sets default filePath when environment variable is not set", func(t *testing.T) {
		if err := os.Unsetenv("YAML_VAULT_FILE_PATH"); err != nil {
			t.Fatalf("Failed to unset environment variable: %v", err)
		}

		client := &YAMLVaultClient{}
		err := client.InitializeClient()
		if err != nil {
			t.Fatalf("Expected no error, but got %v", err)
		}

		if client.filePath != DEFAULT_YAML_FILE_PATH {
			t.Errorf("Expected filePath to be '%s', but got '%s'", DEFAULT_YAML_FILE_PATH, client.filePath)
		}
	})
}

func TestYAMLGetSecrets(t *testing.T) {
	t.Run("returns secrets from YAML file", func(t *testing.T) {
		expectedSecrets := []VaultSecret{
			{Name: "secret1", Value: "value1"},
			{Name: "secret2", Value: "value2"},
		}

		file, err := os.CreateTemp("", "secrets.yaml")
		if err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}
		defer func() {
			if err := os.Remove(file.Name()); err != nil {
				t.Fatalf("Failed to remove temp file: %v", err)
			}
		}()

		bytes, err := yaml.Marshal(expectedSecrets)
		if err != nil {
			t.Fatalf("Failed to marshal secrets: %v", err)
		}

		if _, err := file.Write(bytes); err != nil {
			t.Fatalf("Failed to write to temp file: %v", err)
		}
		defer func() {
			if err := file.Close(); err != nil {
				t.Fatalf("Failed to close temp file: %v", err)
			}
		}()

		client := &YAMLVaultClient{filePath: file.Name()}
		secrets, err := client.GetSecrets()
		if err != nil {
			t.Fatalf("Expected no error, but got %v", err)
		}

		if len(secrets) != len(expectedSecrets) {
			t.Fatalf("Expected %d secrets, but got %d", len(expectedSecrets), len(secrets))
		}

		for i, secret := range secrets {
			if secret.Name != expectedSecrets[i].Name || secret.Value != expectedSecrets[i].Value {
				t.Errorf("Expected secret %v, but got %v", expectedSecrets[i], secret)
			}
		}
	})

	t.Run("returns error when file does not exist", func(t *testing.T) {
		client := &YAMLVaultClient{filePath: "nonexistent.yaml"}
		_, err := client.GetSecrets()
		if err == nil {
			t.Fatalf("Expected error, but got nil")
		}
	})

	t.Run("returns error when YAML is invalid", func(t *testing.T) {
		file, err := os.CreateTemp("", "secrets.yaml")
		if err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}
		defer func() {
			if err := os.Remove(file.Name()); err != nil {
				t.Fatalf("Failed to remove temp file: %v", err)
			}
		}()

		if _, err := file.Write([]byte("invalid yaml")); err != nil {
			t.Fatalf("Failed to write to temp file: %v", err)
		}
		defer func() {
			if err := file.Close(); err != nil {
				t.Fatalf("Failed to close temp file: %v", err)
			}
		}()

		client := &YAMLVaultClient{filePath: file.Name()}
		_, err = client.GetSecrets()
		if err == nil {
			t.Fatalf("Expected error, but got nil")
		}
	})
}
