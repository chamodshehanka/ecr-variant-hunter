package services

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
	"github.com/chamodshehanka/ecr-variant-hunter/internal/config"
	"github.com/sirupsen/logrus"
	"sync"
)

func DeleteOutdatedImages(ecrClient *ecr.Client) {
	var wg sync.WaitGroup
	reposList := config.EnvValues.RepositoryList
	olderThanDays := config.EnvValues.ImagesRetentionDays

	if len(reposList) == 0 {
		logrus.Fatalf("No repositories found in the configuration")
		return
	}

	if olderThanDays == 0 {
		logrus.Fatalf("No retention days found in the configuration")
		return
	}

	for _, repo := range reposList {
		logrus.Infof("Deleting outdated images for repository: %s", repo)
		wg.Add(1)
		go func(repositoryName string) {
			defer wg.Done()
			ctx := context.TODO()
			DeleteECROldImages(ctx, ecrClient, repositoryName, olderThanDays)
		}(repo)
	}

	wg.Wait()
	logrus.Info("Completed deleting outdated images for all repositories")
}
