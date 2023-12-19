package secretsmanagerp

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

var client *secretsmanager.Client

func ensureClient() error {
	if client == nil {
		cfg, err := config.LoadDefaultConfig(context.Background())

		if err != nil {
			return fmt.Errorf("failed to load AWS config: %w", err)
		}

		client = secretsmanager.NewFromConfig(cfg)
	}

	return nil
}

func getSecretString(secretName string) (secretString string, err error) {
	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretName),
	}

	result, err := client.GetSecretValue(context.Background(), input)

	if err != nil {
		return
	}

	secretString = *result.SecretString
	return
}

func GetSecret(ident string) (value string, err, fatal error) {
	fatal = ensureClient()
	if fatal != nil {
		return
	}

	secretName, field, hasField := strings.Cut(ident, "$")

	if !hasField {
		secretName = ident
	}

	value, err = getSecretString(secretName)

	if err != nil {
		return
	}

	if hasField && field != "" {
		secretJson := map[string]string{}
		err = json.Unmarshal([]byte(value), &secretJson)

		if err != nil {
			return
		}

		var ok bool
		value, ok = secretJson[field]
		if !ok {
			err = fmt.Errorf("field %s does not exist on secret %s", field, secretName)
		}
	}

	return
}
