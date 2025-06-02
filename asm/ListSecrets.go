package asm

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// listSecrets retrieves all secrets from AWS Secrets Manager in the specified region.

func ListSecrets(sess *session.Session, region string, ctx context.Context) ([]string, error) {

	logging := log.FromContext(ctx)

	// Create a new Secrets Manager client
	svc := secretsmanager.New(sess, &aws.Config{Region: aws.String(region)})

	// List all secrets
	input := &secretsmanager.ListSecretsInput{}
	var secrets []string
	err := svc.ListSecretsPages(input, func(page *secretsmanager.ListSecretsOutput, lastPage bool) bool {
		for _, secret := range page.SecretList {
			secrets = append(secrets, aws.StringValue(secret.Name))
		}
		return !lastPage
	})
	if err != nil {
		logging.Error(err, "failed to list secrets")
		return nil, err
	}

	return secrets, nil

}
