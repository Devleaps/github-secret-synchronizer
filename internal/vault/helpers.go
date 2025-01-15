package vault

import "strings"

func FormatSecretName(secretName string) string {
	return strings.ToUpper(strings.ReplaceAll(strings.ReplaceAll(secretName, " ", "_"), "-", "_"))
}
