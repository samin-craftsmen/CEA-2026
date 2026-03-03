package database

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/samin-craftsmen/meal-headcount-planner-backend/config"
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

	// Create AWS session
	sess, err := session.NewSession(&aws.Config{
		Region:   aws.String(cfg.DynamoDBRegion),
		Endpoint: aws.String(cfg.DynamoDBEndpoint), // For local DynamoDB
	})
	if err != nil {
		return fmt.Errorf("failed to create AWS session: %w", err)
	}

	// Create DynamoDB client
	DBClient = dynamodb.New(sess)

	// Optional: check connection by listing tables
	_, err = DBClient.ListTables(&dynamodb.ListTablesInput{})
	if err != nil {
		return fmt.Errorf("failed to list tables: %w", err)
	}

	fmt.Println("✅ DynamoDB initialized successfully")
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
