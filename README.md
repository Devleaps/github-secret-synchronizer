# github-secret-synchronizer

An application to synchronize secrets/variables into your GitHub organization.

## Features

### Typings
- [x] Synchronize to GitHub as a secret
- [x] Synchronize to GitHub as a variable

### Visibility
- [x] Synchronize to GitHub, visible to all repositories
- [x] Synchronize to GitHub, visible to private repositories
- [x] Synchronize to GitHub, visible to selected repositories

### Defaults
- [x] Set a default visibility (all, private, selected)
- [x] Set a default typing (secret, variable)
- [x] Set a default (comma-separated) string of repositories

### Vaults
- [x] Pull information from a JSON file to synchronize
- [x] Pull information from a YAML file to synchronize
- [x] Pull information from Azure Key Vault to synchronize
- [x] Pull information from AWS Secrets Manager to synchronize

## To be implemented

### Vaults
- [ ] Pull information from HashiCorp Vault to synchronize
- [ ] Pull information from Keeper Security to synchronize

### Synchronizing
- [ ] Synchronize by polling a source
- [ ] Schedule synchronizations using cron
- [ ] Allow for a dry-run
- [ ] Allow for full management (e.g. deleting secrets that are not in the source)

## GitHub

For github-secrets-synchronizer to work, it needs to be able to write secrets/variables to your GitHub organization. This can be done by creating a GitHub App and installing it on your organization. GitHub has excellent documentation on how to do this: [Creating a GitHub App](https://docs.github.com/en/apps/creating-github-apps).

Whilst creating the GitHub App, make sure to give it the following permissions:
- `org:secrets:write` - if you want to synchronize secrets
- `org:variables:write` - if you want to synchronize variables
- `repository:metadata:read` - if you want to synchronize to selected repositories

Make sure to install the GitHub App on the repositories you want to synchronize to in case you want to synchronize to selected repositories.

The following environment variables are mandatory for github-secrets-synchronizer to work:
- `GITHUB_APP_ID` - The ID of the GitHub App
- `GITHUB_APP_PRIVATE_KEY` - The private key of the GitHub App
- `GITHUB_INSTALLATION_ID` - The installation ID of the GitHub App
- `GITHUB_ORGANIZATION` - The organization to synchronize to

## Defaults

Defaults are naturally set within github-secrets-synchronizer. These defaults can be overriden by the user by setting the following environment variables:
- `DEFAULT_VISIBILITY` - The default visibility of the secret/variable (all, private, selected)
- `DEFAULT_TYPE` - The default type of the secret/variable (secret, variable)
- `DEFAULT_REPOSITORIES` - The default repositories that the secret/variable should be visible to (comma-separated)

The default visibility is `all`, the default type is `secret` and the default repositories is an empty string.

## Vaults

To choose a vault to use, set the `VAULT_TYPE` environment variable to either `json`, `yaml`, `azure` or `aws`. Depending on the vault you choose, you need to set the necessary environment variables.

### JSON

Secrets/Variables can be synchronized from a JSON file. For github-secrets-synchronizer to know where this file is, you can pass the following environment variable:
- `JSON_VAULT_FILE_PATH` - The path to the JSON file

The default is `secrets.json`.

This JSON file should be an array of objects, where each object has the following properties:
- `name` (string) **required**: The name of the secret/variable
- `value` (string) **required**: The value of the secret/variable
- `visibility` (string) **optional**: The visibility of the secret/variable (all, private, selected)
- `type` (string) **optional**: The type of the secret/variable (secret, variable)
- `repositories` (string) **optional**: The repositories that the secret/variable should be visible to (comma-separated)

```json
[
    {
        "name": "SECRET_ONE",
        "value": "some-random-value",
        "visibility": "all",
        "type": "secret"
    },
    {
        "name": "VARIABLE_ONE",
        "value": "some-random-value",
        "visibility": "private",
        "type": "variable"
    },
    {
        "name": "SECRET_TWO",
        "value": "some-random-value",
        "visibility": "selected",
        "type": "secret",
        "repositories": [
            "repo-one",
            "repo-two"
        ]
    },
    {
        "name": "DEFAULT_ONE",
        "value": "some-random-value",
    }
]
```

An example with more possibilities can be found in [`examples/secrets.json`](examples/secrets.json).


### YAML

Secrets/Variables can be synchronized from a YAML file. For github-secrets-synchronizer to know where this file is, you can pass the following environment variable:
- `YAML_VAULT_FILE_PATH` - The path to the YAML file

The default is `secrets.yaml`.

This YAML file should be an array of objects, where each object has the following properties:
- `name` (string) **required**: The name of the secret/variable
- `value` (string) **required**: The value of the secret/variable
- `visibility` (string) **optional**: The visibility of the secret/variable (all, private, selected)
- `type` (string) **optional**: The type of the secret/variable (secret, variable)
- `repositories` (string) **optional**: The repositories that the secret/variable should be visible to (comma-separated)

```yaml
- name: SECRET_ONE
  value: some-random-value
  visibility: all
  type: secret
- name: VARIABLE_ONE
  value: some-random-value
  visibility: private
  type: variable
- name: SECRET_TWO
  value: some-random-value
  visibility: selected
  type: secret
  repositories:
    - repo-one
    - repo-two
- name: DEFAULT_ONE
  value: some-random-value
```

An example with more possibilities can be found in [`examples/secrets.yaml`](examples/secrets.yaml).

### Azure Key Vault

Secrets stored in Azure Key Vault can be synchronized towards GitHub Secrets/Variables. The following environment variables are mandatory for github-secrets-synchronizer to work with Azure Key Vault:
- `AZURE_CLIENT_ID` - The client ID of the Azure Service Principal
- `AZURE_CLIENT_SECRET` - The client secret of the Azure Service Principal
- `AZURE_TENANT_ID` - The tenant ID of the Azure Service Principal
- `AZURE_KEYVAULT_URL` - The URL of the Azure Key Vault

Within Azure, you can create an App Registration and assign it the necessary permissions to access the Key Vault. Once you created an App Registration, you can modify the Key Vault's access policies to allow the App Registration to read secrets. The least permissive role you can assign is `Key Vault Secrets User`.

Once you have the necessary permissions, you can add secrets to the Key Vault. To set the visibility of the secret, you can add a tag to the secret with the key `visibility` and the value `all`, `private` or `selected`. To set the type of the secret, you can add a tag to the secret with the key `type` and the value `secret` or `variable`. To set the repositories the secret should be visible to, you can add a tag to the secret with the key `repositories` and the value of the repositories (comma-separated).

### AWS Secrets Manager

Secrets stored in AWS Secrets Manager can be synchronized towards GitHub Secrets/Variables. The following environment variables are mandatory for github-secrets-synchronizer to work with AWS Secrets Manager:
- `AWS_ACCESS_KEY_ID` - The access key ID of the AWS IAM user
- `AWS_SECRET_ACCESS_KEY` - The secret access key of the AWS IAM user
- `AWS_REGION` - The region of the AWS Secrets Manager

Within AWS, you can create an IAM user and assign it the necessary permissions to access the secrets within the Secrets Manager. Note that this has only been tested with plain-text secrets. 

Once you have the necessary permissions, you can add secrets to the Secrets Manager. To set the visibility of the secret, you can add a tag to the secret with the key `visibility` and the value `all`, `private` or `selected`. To set the type of the secret, you can add a tag to the secret with the key `type` and the value `secret` or `variable`. To set the repositories the secret should be visible to, you can add a tag to the secret with the key `repositories` and the value of the repositories (comma-separated).

## Local Development

To run github-secrets-synchronizer locally, you can use the following command:
```bash
go run main.go
```

To set the necessary environment variables to be able to run github-secrets-synchronizer locally, you can create a `.env` file in the root of the project.

By setting an environment variable `LOCAL` to `true`, github-secrets-synchronizer will pretty-print the output to the console.
