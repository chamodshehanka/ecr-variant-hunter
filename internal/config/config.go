package config

type AwsConfig struct {
	Region          string
	AccessKeyID     string
	SecretAccessKey string
	RegistryURL     string
}

type Config struct {
	AWS                 AwsConfig
	RepositoryList      []string
	ImagesRetentionDays int
}
