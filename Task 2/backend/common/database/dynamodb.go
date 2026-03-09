package database

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/samin-craftsmen/meal-headcount-planner-backend/common/config"
)

var (
	// DBClient is the global DynamoDB client
	DBClient *dynamodb.DynamoDB
	// TableName stores the table to be used
	TableName string
)

// InitDynamoDB initializes the DynamoDB client
func InitDynamoDB(cfg *config.Config) error {
	TableName = cfg.DynamoDBTable

	awsConfig := &aws.Config{
		Region: aws.String(cfg.DynamoDBRegion),
	}
	if cfg.DynamoDBEndpoint != "" {
		awsConfig.Endpoint = aws.String(cfg.DynamoDBEndpoint)
	}

	// Create AWS session
	sess, err := session.NewSession(awsConfig)
	if err != nil {
		return fmt.Errorf("failed to create AWS session: %w", err)
	}

	// Create DynamoDB client. initialize only if not already set
	if DBClient == nil {
		DBClient = dynamodb.New(sess)
	}
	fmt.Println("DynamoDB initialized successfully")
	return nil
}

// GetDBClient returns the DynamoDB client
func GetDBClient() *dynamodb.DynamoDB {
	return DBClient
}

// GetTableName returns the configured table name
func GetTableName() string {
	return TableName
}
