package services

import (
	"context"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
	"github.com/aws/aws-sdk-go-v2/service/ecr/types"
	"github.com/chamodshehanka/ecr-variant-hunter/internal/config"
	"github.com/sirupsen/logrus"
	"time"
)

func GetECRConfig() *ecr.Client {
	ctx := context.TODO()
	region := config.EnvValues.AWS.Region
	accessKey := config.EnvValues.AWS.AccessKeyID
	secretKey := config.EnvValues.AWS.SecretAccessKey

	cfg, err := awsConfig.LoadDefaultConfig(ctx,
		awsConfig.WithRegion(region),
		awsConfig.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(accessKey, secretKey, ""),
		),
	)

	if err != nil {
		logrus.Fatalf("Error loading config: %v", err)
		return nil
	}

	return ecr.NewFromConfig(cfg)
}

func DeleteECROldImages(ctx context.Context, client *ecr.Client, repositoryName string, olderThanDays int) {
	// List all images in the repository
	listImagesInput := &ecr.ListImagesInput{
		RepositoryName: &repositoryName,
	}

	listImagesOutput, err := client.ListImages(ctx, listImagesInput)
	if err != nil {
		logrus.Errorf("Error listing images: %v", err)
		return
	}

	// Delete images older than the specified number of days
	for _, imageID := range listImagesOutput.ImageIds {
		describeImagesInput := &ecr.DescribeImagesInput{
			RepositoryName: &repositoryName,
			ImageIds:       []types.ImageIdentifier{imageID},
		}
		describeImagesOutput, err := client.DescribeImages(ctx, describeImagesInput)
		if err != nil {
			logrus.Errorf("Error describing image: %v", err)
			return
		}

		imagePushedAt := describeImagesOutput.ImageDetails[0].ImagePushedAt
		if imagePushedAt.Before(time.Now().AddDate(0, 0, -olderThanDays)) {
			logrus.Infof("Deleting image: %v", *imageID.ImageTag)

			deleteImageInput := &ecr.BatchDeleteImageInput{
				RepositoryName: &repositoryName,
				ImageIds:       []types.ImageIdentifier{imageID},
			}

			_, err := client.BatchDeleteImage(ctx, deleteImageInput)
			if err != nil {
				logrus.Errorf("Error deleting image: %v", err)
				return
			}
		}
	}
}
