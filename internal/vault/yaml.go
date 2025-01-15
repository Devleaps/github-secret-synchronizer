package vault

import (
	"fmt"
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

const (
	DEFAULT_YAML_FILE_PATH = "secrets.yaml"
)

// YAMLVaultClient implements the VaultClient interface for a local YAML file
type YAMLVaultClient struct {
	filePath string
}

// InitializeClient initializes the local vault client by reading the YAML file
func (v *YAMLVaultClient) InitializeClient() error {
	filePath := os.Getenv("YAML_VAULT_FILE_PATH")

	if filePath == "" {
		filePath = DEFAULT_YAML_FILE_PATH
	}

	v.filePath = filePath

	return nil
}

// GetSecrets retrieves secrets from the local YAML file
func (v *YAMLVaultClient) GetSecrets() ([]VaultSecret, error) {
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

	if err := yaml.Unmarshal(bytes, &secrets); err != nil {
		return nil, fmt.Errorf("failed to unmarshal YAML: %w", err)
	}

	return secrets, nil
}
