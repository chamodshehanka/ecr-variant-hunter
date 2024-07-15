package services

import (
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/chamodshehanka/ecr-variant-hunter/internal/config"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

func DeleteOutdatedImages() {
	var wg sync.WaitGroup

	ecrService, err := GetECRService()
	if err != nil {
		logrus.Fatalf("Error creating ECR service: %s", err)
		return
	}

	for _, repo := range config.EnvValues.RepositoryList {
		logrus.Infof("Deleting old images from repository %s", repo)
		wg.Add(1)

		repoCopy := repo // Create a copy of repo for the goroutine
		go func(repoName string) {
			defer wg.Done()
			deleteOldImages(ecrService, &repoName)
		}(repoCopy)
	}

	wg.Wait() // Wait for all goroutines to finish
}

func deleteOldImages(svc *ecr.ECR, repoName *string) {
	logrus.Infof("Deleting old images from repository %s", *repoName)
	now := time.Now()

	imageRetentionDays := config.EnvValues.ImagesRetentionDays
	if imageRetentionDays == 0 {
		imageRetentionDays = 14 // Set default value to 14 if not specified
	}
	cutoffDate := now.AddDate(0, 0, -imageRetentionDays)

	listResp, err := svc.ListImages(&ecr.ListImagesInput{
		RepositoryName: repoName,
	})
	if err != nil {
		logrus.Errorf("Error listing images for repository %s: %s\n", *repoName, err)
		return
	}

	var oldImageIds []*ecr.ImageIdentifier
	for _, imageId := range listResp.ImageIds {
		descResp, err := svc.DescribeImages(&ecr.DescribeImagesInput{
			RepositoryName: repoName,
			ImageIds:       []*ecr.ImageIdentifier{imageId},
		})
		if err != nil {
			logrus.Errorf("Error describing image %s in repository %s: %s\n", *imageId.ImageDigest, *repoName, err)
			continue
		}

		for _, imageDetail := range descResp.ImageDetails {
			if imageDetail.ImagePushedAt.Before(cutoffDate) {
				oldImageIds = append(oldImageIds, &ecr.ImageIdentifier{
					ImageDigest: imageDetail.ImageDigest,
				})
			}
		}
	}

	if len(oldImageIds) > 0 {
		_, err = svc.BatchDeleteImage(&ecr.BatchDeleteImageInput{
			RepositoryName: repoName,
			ImageIds:       oldImageIds,
		})
		if err != nil {
			logrus.Errorf("Error deleting images from repository %s: %s\n", *repoName, err)
			return
		}
		logrus.Infof("Deleted %d images from repository %s\n", len(oldImageIds), *repoName)
	}
}
