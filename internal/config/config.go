package config

type AwsConfig struct {
	Region          string
	AccessKeyID     string
	SecretAccessKey string
}

type Config struct {
	AWS                 AwsConfig
	RepositoryList      []string
	ImagesRetentionDays int
}
