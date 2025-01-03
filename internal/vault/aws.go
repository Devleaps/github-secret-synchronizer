package vault

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager/types"
	"github.com/rs/zerolog/log"
)

// AWSVaultClient implements the VaultClient interface for an AWS Vault
type AWSVaultClient struct {
	client *secretsmanager.Client
}

// InitializeClient initializes the AWS vault client
func (v *AWSVaultClient) InitializeClient() error {
	cfg, err := config.LoadDefaultConfig(context.Background())

	if err != nil {
		return err
	}

	v.client = secretsmanager.NewFromConfig(cfg)

	return err
}

// GetSecrets retrieves secrets from the AWS Vault instance
func (v *AWSVaultClient) GetSecrets() ([]VaultSecret, error) {
	var secrets []VaultSecret
	// TODO: Fix pagination
	awsSecrets, err := v.client.ListSecrets(context.Background(), &secretsmanager.ListSecretsInput{})

	if err != nil {
		return nil, err
	}

	if awsSecrets.NextToken != nil {
		log.Warn().Msg("AWS Vault secrets list is truncated")
	}

	for _, secret := range awsSecrets.SecretList {
		repositories := v.retrieveTag(secret.Tags, "repositories")
		visibility := v.retrieveTag(secret.Tags, "visibility")
		secretType := v.retrieveTag(secret.Tags, "type")

		// TODO: Implement batch getting of secrets
		result, err := v.client.GetSecretValue(context.Background(), &secretsmanager.GetSecretValueInput{
			SecretId: secret.Name,
		})

		if err != nil {
			return nil, err
		}

		secrets = append(secrets, VaultSecret{
			Name:         strings.ToUpper(strings.ReplaceAll(strings.ReplaceAll(*secret.Name, " ", "_"), "-", "_")),
			Value:        *result.SecretString,
			Type:         secretType,
			Visibility:   visibility,
			Repositories: strings.Split(repositories, ","),
		})
	}

	return secrets, nil
}

func (v *AWSVaultClient) retrieveTag(tags []types.Tag, key string) string {
	for _, tag := range tags {
		if *tag.Key == key {
			return *tag.Value
		}
	}
	return ""
}
