package config

import "os"

type Config struct {
	Port           string
	AWSRegion      string
	DynamoEndpoint string
	TableName      string
	JWTSecret      string
}

func Load() *Config {
	return &Config{
		Port:           getEnv("PORT", ":8080"),
		AWSRegion:      getEnv("AWS_REGION", "us-east-1"),
		DynamoEndpoint: getEnv("DYNAMO_ENDPOINT", "http://dynamodb-local:8000"),
		TableName:      getEnv("TABLE_NAME", "LinkTable"),
		JWTSecret:      getEnv("JWT_SECRET", "dev-secret-change-in-prod"),
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}