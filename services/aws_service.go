package services

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/chamodshehanka/ecr-variant-hunter/internal/config"
)

func GetECRService() (*ecr.ECR, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: &config.EnvValues.AWS.Region,
		Credentials: credentials.NewStaticCredentials(
			config.EnvValues.AWS.AccessKeyID,
			config.EnvValues.AWS.SecretAccessKey,
			"",
		),
	})

	if err != nil {
		return nil, err
	}

	return ecr.New(sess), nil
}
