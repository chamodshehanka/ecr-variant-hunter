package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"os"
	"strings"
)

var config Config

var requiredEnvs = [...]string{
	"AWS_REGION",
	"AWS_ACCESS_KEY_ID",
	"AWS_SECRET_ACCESS_KEY",
	"ECR_REPOS_LIST",
}

func LoadConfig() error {
	config = Config{
		AWS: AwsConfig{
			Region:          os.Getenv("AWS_REGION"),
			AccessKeyID:     os.Getenv("AWS_ACCESS_KEY_ID"),
			SecretAccessKey: os.Getenv("AWS_SECRET_ACCESS_KEY"),
			RegistryURL:     os.Getenv("AWS_REGISTRY_URL"),
		},
	}

	repoList := os.Getenv("ECR_REPOS_LIST")
	if repoList != "" {
		repos := strings.Split(repoList, ",")

		for i, repo := range repos {
			repos[i] = fmt.Sprintf("%s/%s", config.AWS.RegistryURL, repo)
		}

		config.RepositoryList = repos
	}

	logrus.Infoln("Config loaded successfully!")
	return nil
}

func ensureRequiredEnvsAreAvailable() error {
	err := godotenv.Load()
	if err != nil {
		logrus.Fatalf("Error loading .env file")
	}

	for _, env := range requiredEnvs {
		if getEnv(env) == "" {
			return fmt.Errorf("fatal: required environment variable '%s' not found", env)
		}
	}
	return nil
}

func getEnv(key string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return ""
}
