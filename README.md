# github-secret-synchronizer

An application to synchronize secrets to a GitHub organization.

## Functionality

### Typings
- [x] Able to synchronize secrets to a GitHub organization
- [ ] Able to synchronize variables to a GitHub organization

### Visibility
- [x] Able to synchronize secrets to a GitHub organization with `all` visibility
- [ ] Able to synchronize secrets to a GitHub organization with `private` visibility
- [ ] Able to synchronize secrets to a GitHub organization with `selected` visibility

### Vaults
- [x] Pull information from a JSON file to synchronize
- [x] Pull information from a YAML file to synchronize
- [ ] Pull information from Azure Key Vault to synchronize
- [ ] Pull information from HashiCorp Vault to synchronize
- [ ] Pull information from AWS Secrets Manager to synchronize
- [ ] Pull information from Keeper Security to synchronize

### Defaults
- [ ] Set a default visibility (all, private, selected)
- [ ] Set a default typing (secret, variable)

### Synchronizing
- [ ] Allow for polling to synchronize
- [ ] Allow for cron to synchronize
- [ ] Allow for a dry-run
- [ ] Allow for full management (e.g. deleting secrets that are not in the source)
