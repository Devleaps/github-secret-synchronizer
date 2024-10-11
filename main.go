package main

import (
	"log"

	"github.com/Devleaps/github-secret-synchronizer/internal/github"
	"github.com/Devleaps/github-secret-synchronizer/internal/vault"
	"github.com/joho/godotenv"
)

func main() {

	// Use godotenv to load the .env file
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	githubWrapper := &github.GitHubWrapper{}

	err := githubWrapper.NewClient()
	if err != nil {
		log.Fatalf("Error initializing GitHub client: %v", err)
	}

	// Initialize the vault client
	vaultClient, err := vault.InitializeVaultClient()
	if err != nil {
		log.Fatalf("Error initializing vault client: %v", err)
	}

	vaultClient.InitializeClient()

	// Retrieve secrets from the vault
	secrets, err := vaultClient.GetSecrets()
	if err != nil {
		log.Fatalf("Error retrieving secrets from vault: %v", err)
	}

	// Encrypt and push each secret to GitHub
	for _, secret := range secrets {
		log.Printf("Pushing secret %s...", secret.Name)
		err = githubWrapper.PushSecret(secret.Name, secret.Value, secret.Visibility)
		if err != nil {
			log.Fatalf("Error pushing encrypted secret: %v", err)
		}

		log.Printf("Secret %s pushed successfully", secret.Name)
	}
}
