package main

import (
	"log"

	"github.com/Devleaps/github-secrets-synchronizer/internal/github"
	"github.com/Devleaps/github-secrets-synchronizer/internal/vault/local"
	"github.com/joho/godotenv"
)

func main() {

	// Use godotenv to load the .env file
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	githubService := github.NewGitHubService()

	client, err := githubService.NewClient()
	if err != nil {
		log.Fatalf("Error initializing GitHub client: %v", err)
	}

	publicKey, keyID, err := githubService.GetOrgPublicKey(client)
	if err != nil {
		log.Fatalf("Error retrieving organization public key: %v", err)
	}

	vaultClient := &local.LocalVaultClient{}
	// Initialize the vault client
	if err := vaultClient.InitializeClient(); err != nil {
		log.Fatalf("Error initializing vault client: %v", err)
	}

	// Retrieve secrets from the vault
	secrets, err := vaultClient.GetSecrets()
	if err != nil {
		log.Fatalf("Error retrieving secrets from vault: %v", err)
	}

	// Encrypt and push each secret to GitHub
	for _, secret := range *secrets {
		log.Printf("Pushing secret %s...", secret.Name)
		encryptedSecret, err := githubService.EncryptSecret(secret.Value, publicKey)
		if err != nil {
			log.Fatalf("Error encrypting secret: %v", err)
		}

		err = githubService.PushEncryptedSecret(client, secret.Name, encryptedSecret, keyID, secret.Visibility)
		if err != nil {
			log.Fatalf("Error pushing encrypted secret: %v", err)
		}

		log.Printf("Secret %s pushed successfully", secret.Name)
	}
}
