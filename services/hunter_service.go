package services

import (
	"fmt"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/chamodshehanka/ecr-variant-hunter/internal/config"
	"github.com/sirupsen/logrus"
	"time"
)

func DeleteOutdatedImages() {
	//var wg sync.WaitGroup

	// Example: List all repositories
	for _, repo := range config.EnvValues.RepositoryList {

		logrus.Infof("Deleting old images from repository %s", repo)
		//wg.Add(1)
		//
		//go func(repoName *string) {
		//	//defer wg.Done()
		//	deleteOldImages(svc, repoName)
		//}(&repo)
	}
}

func deleteOldImages(svc *ecr.ECR, repoName *string) {
	logrus.Infof("Deleting old images from repository %s", *repoName)
	now := time.Now()
	cutoffDate := now.AddDate(0, 0, -14) // 14 days ago

	// List images
	listResp, err := svc.ListImages(&ecr.ListImagesInput{
		RepositoryName: repoName,
	})
	if err != nil {
		fmt.Printf("Error listing images for repository %s: %s\n", *repoName, err)
		return
	}

	var oldImageIds []*ecr.ImageIdentifier
	for _, imageId := range listResp.ImageIds {
		descResp, err := svc.DescribeImages(&ecr.DescribeImagesInput{
			RepositoryName: repoName,
			ImageIds:       []*ecr.ImageIdentifier{imageId},
		})
		if err != nil {
			fmt.Printf("Error describing image %s in repository %s: %s\n", *imageId.ImageDigest, *repoName, err)
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
			fmt.Printf("Error deleting images from repository %s: %s\n", *repoName, err)
			return
		}
		fmt.Printf("Deleted %d images from repository %s\n", len(oldImageIds), *repoName)
	}
}
