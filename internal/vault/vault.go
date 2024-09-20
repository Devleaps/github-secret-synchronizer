package vault

const (
	ALL_VISIBILITY     = "all"
	PRIVATE_VISIBILITY = "private"
)

// VaultClient defines the interface for interacting with different vaults
type VaultClient interface {
	InitializeClient() error
	GetSecrets() (*[]VaultSecret, error)
}

type VaultSecret struct {
	Name       string `json:"name"`
	Value      string `json:"value"`
	Visibility string `json:"visibility"`
}
