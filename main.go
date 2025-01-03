package main

import (
	"os"

	"github.com/Devleaps/github-secret-synchronizer/internal/github"
	"github.com/Devleaps/github-secret-synchronizer/internal/vault"
	_ "github.com/joho/godotenv/autoload"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func init() {
	if os.Getenv("LOCAL") == "true" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}
}

func main() {
	githubWrapper := &github.GitHubWrapper{}

	err := githubWrapper.NewClient()
	if err != nil {
		log.Fatal().Err(err).Msg("Error initializing GitHub client")
	}

	// Initialize the vault
	vaultClient, err := vault.InitializeVault()
	if err != nil {
		log.Fatal().Err(err).Msg("Error initializing vault client")
	}

	vaultClient.InitializeClient()

	// Retrieve secrets from the vault
	secrets, err := vaultClient.GetSecrets()
	if err != nil {
		log.Fatal().Err(err).Msg("Error retrieving secrets from vault")
	}

	// Set defaults in case they are not provided
	vault.HandleDefaults(secrets)

	// Encrypt and push each secret to GitHub
	for _, secret := range secrets {
		if secret.Type == vault.SECRET_TYPE {
			log.Info().Str("secret", secret.Name).Msg("Pushing secret...")
			err = githubWrapper.PushSecret(secret)
			if err != nil {
				log.Fatal().Err(err).Msg("Error pushing secret")
			}
			log.Info().Str("secret", secret.Name).Msg("Secret pushed successfully!")
		} else if secret.Type == vault.VARIABLE_TYPE {
			log.Info().Str("variable", secret.Name).Msg("Pushing variable...")
			err = githubWrapper.PushVariable(secret)
			if err != nil {
				log.Fatal().Err(err).Msg("Error pushing variable")
			}
			log.Info().Str("variable", secret.Name).Msg("Variable pushed successfully!")
		}
	}
}
