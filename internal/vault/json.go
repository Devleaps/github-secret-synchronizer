package vault

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

const (
	DEFAULT_JSON_FILE_PATH = "secrets.json"
)

// JSONVaultClient implements the VaultClient interface for a local JSON file
type JSONVaultClient struct {
	filePath string
}

// InitializeClient initializes the local vault client by reading the JSON file
func (v *JSONVaultClient) InitializeClient() error {
	filePath := os.Getenv("JSON_VAULT_FILE_PATH")

	if filePath == "" {
		filePath = DEFAULT_JSON_FILE_PATH
	}

	v.filePath = filePath

	return nil
}

// GetSecrets retrieves secrets from the local JSON file
func (v *JSONVaultClient) GetSecrets() ([]VaultSecret, error) {
	file, err := os.Open(v.filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var secrets []VaultSecret

	if err := json.Unmarshal(bytes, &secrets); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return secrets, nil
}
