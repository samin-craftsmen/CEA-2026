package repository

import (
	"fmt"
	"strings"

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

	const defaultTeamID = "TEAM#engineering"

	_, err = db.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String(table),
		Item: map[string]*dynamodb.AttributeValue{
			"PK":     {S: aws.String(pk)},
			"SK":     {S: aws.String("META")},
			"role":   {S: aws.String("EMPLOYEE")},
			"teamId": {S: aws.String(defaultTeamID)},
		},
	})
	if err != nil {
		return err
	}

	// Also create the team membership record so VerifyTeamMembership works.
	_, err = db.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String(table),
		Item: map[string]*dynamodb.AttributeValue{
			"PK": {S: aws.String(defaultTeamID)},
			"SK": {S: aws.String("USER#" + discordID)},
		},
	})
	return err
}

func GetUserMeals(userID string, date string) (map[string]string, error) {
	db := database.GetDBClient()
	table := database.GetTableName()

	pk := "DAY#" + date
	userFilter := "USER#" + userID

	out, err := db.Query(&dynamodb.QueryInput{
		TableName:              aws.String(table),
		KeyConditionExpression: aws.String("PK = :pk"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":pk": {S: aws.String(pk)},
		},
	})
	if err != nil {
		return nil, err
	}

	result := map[string]string{}
	for _, item := range out.Items {
		sk := *item["SK"].S
		if !strings.Contains(sk, userFilter) {
			continue
		}
		// SK format: TEAM#<teamId>#MEAL#<mealType>#USER#<userId>
		mealIdx := strings.Index(sk, "#MEAL#")
		userIdx := strings.Index(sk, "#USER#")
		if mealIdx >= 0 && userIdx > mealIdx {
			mealType := sk[mealIdx+6 : userIdx]
			result[mealType] = *item["participation"].S
		}
	}
	return result, nil
}

func GetUserRole(discordID string) (string, string, error) {
	db := database.GetDBClient()
	table := database.GetTableName()

	result, err := db.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(table),
		Key: map[string]*dynamodb.AttributeValue{
			"PK": {S: aws.String("USER#" + discordID)},
			"SK": {S: aws.String("META")},
		},
	})
	if err != nil {
		return "", "", err
	}
	if result.Item == nil {
		return "", "", fmt.Errorf("user not found")
	}
	role := ""
	if v, ok := result.Item["role"]; ok && v.S != nil {
		role = *v.S
	}
	teamID := ""
	if v, ok := result.Item["teamId"]; ok && v.S != nil {
		teamID = *v.S
	}
	if role == "" {
		return "", "", fmt.Errorf("role not found for user")
	}
	if teamID == "" {
		return "", "", fmt.Errorf("teamId not found for user")
	}
	return role, teamID, nil
}

func SetMealParticipation(userID, date, mealType, status, teamID string) error {
	db := database.GetDBClient()
	table := database.GetTableName()

	sk := teamID + "#MEAL#" + mealType + "#USER#" + userID

	_, err := db.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String(table),
		Item: map[string]*dynamodb.AttributeValue{
			"PK":            {S: aws.String("DAY#" + date)},
			"SK":            {S: aws.String(sk)},
			"participation": {S: aws.String(status)},
		},
	})
	return err
}

// VerifyTeamMembership checks whether targetUserID is a member of the given team.
// teamID must be in the stored format (e.g. "TEAM#engineering").
func VerifyTeamMembership(teamID, targetUserID string) (bool, error) {
	db := database.GetDBClient()
	table := database.GetTableName()

	result, err := db.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(table),
		Key: map[string]*dynamodb.AttributeValue{
			"PK": {S: aws.String(teamID)},
			"SK": {S: aws.String("USER#" + targetUserID)},
		},
	})
	if err != nil {
		return false, err
	}
	return result.Item != nil, nil
}
