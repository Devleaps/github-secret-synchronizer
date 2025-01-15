package main

import (
	"os"

	"github.com/Devleaps/github-secret-synchronizer/internal/github"
	"github.com/Devleaps/github-secret-synchronizer/internal/vault"
	_ "github.com/joho/godotenv/autoload"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	_ "golang.org/x/crypto/x509roots/fallback"
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
		log.Fatal().Err(err).Msg("Error initializing vault")
	}

	err = vaultClient.InitializeClient()

	if err != nil {
		log.Fatal().Err(err).Msg("Error initializing vault client")
	}

	// Retrieve secrets from the vault
	secrets, err := vaultClient.GetSecrets()
	if err != nil {
		log.Fatal().Err(err).Msg("Error retrieving secrets from vault")
	}

	// Set defaults in case they are not provided
	err = vault.HandleDefaults(secrets)

	if err != nil {
		log.Fatal().Err(err).Msg("Error handling defaults")
	}

	dryRun := false
	// If DRY_RUN is set to true, set a variable
	if os.Getenv("DRY_RUN") == "true" {
		dryRun = true
	}

	// Encrypt and push each secret to GitHub
	for _, secret := range secrets {
		if secret.Type == vault.SECRET_TYPE {
			log.Info().Str("secret", secret.Name).Bool("dry-run", dryRun).Msg("Pushing secret...")
			if !dryRun {
				err = githubWrapper.PushSecret(secret)
				if err != nil {
					log.Fatal().Err(err).Msg("Error pushing secret")
				}
			}
			log.Info().Str("secret", secret.Name).Bool("dry-run", dryRun).Msg("Secret pushed successfully!")
		} else if secret.Type == vault.VARIABLE_TYPE {
			log.Info().Str("variable", secret.Name).Bool("dry-run", dryRun).Msg("Pushing variable...")
			if !dryRun {
				err = githubWrapper.PushVariable(secret)
				if err != nil {
					log.Fatal().Err(err).Msg("Error pushing variable")
				}
			}
			log.Info().Str("variable", secret.Name).Bool("dry-run", dryRun).Msg("Variable pushed successfully!")
		}
	}
}
