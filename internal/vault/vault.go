package vault

import (
	"errors"
	"os"
)

const (
	ALL_VISIBILITY     = "all"
	PRIVATE_VISIBILITY = "private"
)

// VaultClient defines the interface for interacting with different vaults
type VaultClient interface {
	InitializeClient() error
	GetSecrets() ([]VaultSecret, error)
}

type VaultSecret struct {
	Name       string `json:"name"`
	Value      string `json:"value"`
	Visibility string `json:"visibility"`
}

func InitializeVaultClient() (VaultClient, error) {
	vaultType := os.Getenv("VAULT_TYPE")
	switch vaultType {
	case "json":
		return &JSONVaultClient{}, nil
	case "yaml":
		return &YAMLVaultClient{}, nil
	default:
		return nil, errors.New("VAULT_TYPE environment variable is required and must be a valid vault type")
	}
}
