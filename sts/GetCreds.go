package sts

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// GetCredsViaSts retrieves AWS credentials by assuming a specified role in a given account.

func GetCredsViaSts(accessKeyID string, secretAccessKey string, accountId string, roleToAssume string, ctx context.Context) (*session.Session, error) {

	logging := log.FromContext(ctx)

	// build complete role ARN
	roleArn := fmt.Sprintf("arn:aws:iam::%s:role/%s", accountId, roleToAssume)

	// Create a new AWS session
	sess, err := session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(accessKeyID, secretAccessKey, ""),
		Region:      aws.String("eu-west-1"),
	})
	if err != nil {
		logging.Error(err, "failed to create AWS session")
		return nil, err
	}

	// Assume the specified role
	stsSvc := sts.New(sess)
	assumeRoleOutput, err := stsSvc.AssumeRole(&sts.AssumeRoleInput{
		RoleArn:         aws.String(roleArn),
		RoleSessionName: aws.String("ekswatch-session"),
	})
	if err != nil {
		logging.Error(err, "failed to assume role")
		return nil, err
	}

	// Create a new session with the assumed role credentials
	sess, err = session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(
			aws.StringValue(assumeRoleOutput.Credentials.AccessKeyId),
			aws.StringValue(assumeRoleOutput.Credentials.SecretAccessKey),
			aws.StringValue(assumeRoleOutput.Credentials.SessionToken),
		),
		Region: aws.String("eu-west-1"),
	})
	if err != nil {
		logging.Error(err, "failed to create session with assumed role")
		return nil, err
	}
	return sess, nil
}
