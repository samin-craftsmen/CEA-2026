package repository

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/samin-craftsmen/meal-headcount-planner-backend/common/database"
)

func EnsureUserExists(discordID string) error {
	db := database.GetDBClient()
	table := database.GetTableName()
	pk := "USER#" + discordID

	result, err := db.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(table),
		Key: map[string]*dynamodb.AttributeValue{
			"PK": {S: aws.String(pk)},
			"SK": {S: aws.String("META")},
		},
	})
	if err != nil {
		return err
	}
	if result.Item != nil {
		return nil
	}

	_, err = db.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String(table),
		Item: map[string]*dynamodb.AttributeValue{
			"PK":     {S: aws.String(pk)},
			"SK":     {S: aws.String("META")},
			"role":   {S: aws.String("EMPLOYEE")},
			"teamId": {S: aws.String("TEAM#engineering")},
		},
	})
	return err
}

func GetUserMeals(userID string, date string) (map[string]string, error) {
	db := database.GetDBClient()
	table := database.GetTableName()

	pk := "MEAL#" + date
	prefix := "USER#" + userID

	out, err := db.Query(&dynamodb.QueryInput{
		TableName:              aws.String(table),
		KeyConditionExpression: aws.String("PK = :pk AND begins_with(SK, :sk)"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":pk": {S: aws.String(pk)},
			":sk": {S: aws.String(prefix)},
		},
	})
	if err != nil {
		return nil, err
	}

	result := map[string]string{}
	for _, item := range out.Items {
		sk := *item["SK"].S
		if len(sk) >= 5 && sk[len(sk)-5:] == "lunch" {
			result["lunch"] = *item["participation"].S
		}
		if len(sk) >= 6 && sk[len(sk)-6:] == "snacks" {
			result["snacks"] = *item["participation"].S
		}
	}
	return result, nil
}
