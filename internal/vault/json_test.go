package vault

import (
	"encoding/json"
	"os"
	"testing"
)

func TestJSONInitializeClient(t *testing.T) {
	t.Run("sets filePath from environment variable", func(t *testing.T) {
		expectedPath := "custom_secrets.json"
		os.Setenv("JSON_VAULT_FILE_PATH", expectedPath)
		defer os.Unsetenv("JSON_VAULT_FILE_PATH")

		client := &JSONVaultClient{}
		err := client.InitializeClient()
		if err != nil {
			t.Fatalf("Expected no error, but got %v", err)
		}

		if client.filePath != expectedPath {
			t.Errorf("Expected filePath to be '%s', but got '%s'", expectedPath, client.filePath)
		}
	})

	t.Run("sets default filePath when environment variable is not set", func(t *testing.T) {
		os.Unsetenv("JSON_VAULT_FILE_PATH")

		client := &JSONVaultClient{}
		err := client.InitializeClient()
		if err != nil {
			t.Fatalf("Expected no error, but got %v", err)
		}

		if client.filePath != DEFAULT_JSON_FILE_PATH {
			t.Errorf("Expected filePath to be '%s', but got '%s'", DEFAULT_JSON_FILE_PATH, client.filePath)
		}
	})
}
func TestJSONGetSecrets(t *testing.T) {
	t.Run("returns secrets from JSON file", func(t *testing.T) {
		expectedSecrets := []VaultSecret{
			{Name: "secret1", Value: "value1"},
			{Name: "secret2", Value: "value2"},
		}

		file, err := os.CreateTemp("", "secrets.json")
		if err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}
		defer os.Remove(file.Name())

		bytes, err := json.Marshal(expectedSecrets)
		if err != nil {
			t.Fatalf("Failed to marshal secrets: %v", err)
		}

		if _, err := file.Write(bytes); err != nil {
			t.Fatalf("Failed to write to temp file: %v", err)
		}
		file.Close()

		client := &JSONVaultClient{filePath: file.Name()}
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
		client := &JSONVaultClient{filePath: "nonexistent.json"}
		_, err := client.GetSecrets()
		if err == nil {
			t.Fatalf("Expected error, but got nil")
		}
	})

	t.Run("returns error when JSON is invalid", func(t *testing.T) {
		file, err := os.CreateTemp("", "secrets.json")
		if err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}
		defer os.Remove(file.Name())

		if _, err := file.Write([]byte("invalid json")); err != nil {
			t.Fatalf("Failed to write to temp file: %v", err)
		}
		file.Close()

		client := &JSONVaultClient{filePath: file.Name()}
		_, err = client.GetSecrets()
		if err == nil {
			t.Fatalf("Expected error, but got nil")
		}
	})
}
