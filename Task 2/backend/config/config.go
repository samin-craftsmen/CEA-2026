package config

import (
	"os"
)

type Config struct {
	Port             string
	GinMode          string
	DynamoDBEndpoint string
	DynamoDBRegion   string
	DynamoDBTable    string
}

func LoadConfig() *Config {
	return &Config{
		Port:             getEnv("PORT", "8080"),
		GinMode:          getEnv("GIN_MODE", "debug"),
		DynamoDBEndpoint: getEnv("DYNAMODB_ENDPOINT", ""),
		DynamoDBRegion:   getEnv("DYNAMODB_REGION", "us-east-1"),
		DynamoDBTable:    getEnv("DYNAMODB_TABLE", "meal-headcount"),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
