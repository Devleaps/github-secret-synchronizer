package local

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/Devleaps/github-secrets-synchronizer/internal/vault"
)

const (
	DEFAULT_FILE_PATH = "secrets.json"
)

// LocalVaultClient implements the VaultClient interface for a local JSON file
type LocalVaultClient struct {
	filePath string
}

// InitializeClient initializes the local vault client by reading the JSON file
func (v *LocalVaultClient) InitializeClient() error {
	filePath := os.Getenv("LOCAL_VAULT_FILE_PATH")

	if filePath == "" {
		filePath = DEFAULT_FILE_PATH
		log.Printf("Using default file path: %s\n", filePath)
	}

	v.filePath = filePath

	return nil
}

// GetSecrets retrieves secrets from the local JSON file
func (v *LocalVaultClient) GetSecrets() (*[]vault.VaultSecret, error) {
	file, err := os.Open(v.filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	secrets := &[]vault.VaultSecret{}

	if err := json.Unmarshal(bytes, secrets); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return secrets, nil
}
