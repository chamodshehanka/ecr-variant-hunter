package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"os"
	"strconv"
	"strings"
)

var EnvValues Config
var requiredEnvs = [...]string{
	"AWS_REGION",
	"AWS_ACCESS_KEY_ID",
	"AWS_SECRET_ACCESS_KEY",
	"ECR_REPOS_LIST",
}

func LoadConfig() error {
	if _, err := os.Stat(".env"); err == nil {
		err := godotenv.Load()
		if err != nil {
			logrus.Fatalf("Error loading .env file")
		}
	} else if !os.IsNotExist(err) {
		logrus.Fatalf("Error checking .env file: %v", err)
	}

	if err := ensureRequiredEnvsAreAvailable(); err != nil {
		logrus.Fatalf("Error ensuring required environment variables: %s", err)
	}

	EnvValues = Config{
		AWS: AwsConfig{
			Region:          getEnv("AWS_REGION"),
			AccessKeyID:     getEnv("AWS_ACCESS_KEY_ID"),
			SecretAccessKey: getEnv("AWS_SECRET_ACCESS_KEY"),
		},
		ImagesRetentionDays: getNumberEnv("IMAGES_RETENTION_DAYS"),
	}

	repoList := os.Getenv("ECR_REPOS_LIST")
	if repoList != "" {
		repos := strings.Split(repoList, ",")

		for i, repo := range repos {
			repos[i] = repo
		}

		EnvValues.RepositoryList = repos
	}

	logrus.Infoln("Config loaded successfully!")
	return nil
}

func ensureRequiredEnvsAreAvailable() error {
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

func getNumberEnv(key string) int {
	envValue := getEnv(key)
	if envValue != "" {
		number, err := strconv.Atoi(envValue)
		if err != nil {
			logrus.Fatalf("Error converting environment variable '%s' value to int: %s", key, err)
		}
		return number
	}

	return 0
}
