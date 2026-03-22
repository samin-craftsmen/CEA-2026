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

// GetTeamMeals returns all meal participation records for a given team and date.
// Result is a map of userID → map of mealType → participation status.
func GetTeamMeals(teamID, date string) (map[string]map[string]string, error) {
	db := database.GetDBClient()
	table := database.GetTableName()

	pk := "DAY#" + date
	teamPrefix := teamID + "#MEAL#"

	out, err := db.Query(&dynamodb.QueryInput{
		TableName:              aws.String(table),
		KeyConditionExpression: aws.String("PK = :pk AND begins_with(SK, :prefix)"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":pk":     {S: aws.String(pk)},
			":prefix": {S: aws.String(teamPrefix)},
		},
	})
	if err != nil {
		return nil, err
	}

	// SK format: TEAM#<teamId>#MEAL#<mealType>#USER#<userId>
	result := map[string]map[string]string{}
	for _, item := range out.Items {
		sk := *item["SK"].S
		mealIdx := strings.Index(sk, "#MEAL#")
		userIdx := strings.Index(sk, "#USER#")
		if mealIdx < 0 || userIdx <= mealIdx {
			continue
		}
		mealType := sk[mealIdx+6 : userIdx]
		userID := sk[userIdx+6:]
		participation := ""
		if v, ok := item["participation"]; ok && v.S != nil {
			participation = *v.S
		}
		if result[userID] == nil {
			result[userID] = map[string]string{}
		}
		result[userID][mealType] = participation
	}
	return result, nil
}

// GetTeamMembers returns all user IDs that belong to the given team.
// teamID must be in the stored format (e.g. "TEAM#engineering").
func GetTeamMembers(teamID string) ([]string, error) {
	db := database.GetDBClient()
	table := database.GetTableName()

	out, err := db.Query(&dynamodb.QueryInput{
		TableName:              aws.String(table),
		KeyConditionExpression: aws.String("PK = :pk"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":pk": {S: aws.String(teamID)},
		},
	})
	if err != nil {
		return nil, err
	}

	// SK format: USER#<userId>
	var members []string
	for _, item := range out.Items {
		sk := *item["SK"].S
		if strings.HasPrefix(sk, "USER#") {
			members = append(members, sk[5:])
		}
	}
	return members, nil
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

// GetMealTypesForDate returns additional meal types configured for a specific date.
func GetMealTypesForDate(date string) ([]string, error) {
	db := database.GetDBClient()
	table := database.GetTableName()

	out, err := db.Query(&dynamodb.QueryInput{
		TableName:              aws.String(table),
		KeyConditionExpression: aws.String("PK = :pk AND begins_with(SK, :prefix)"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":pk":     {S: aws.String("DAY#" + date)},
			":prefix": {S: aws.String("MEALTYPE#")},
		},
	})
	if err != nil {
		return nil, err
	}

	var types []string
	for _, item := range out.Items {
		sk := *item["SK"].S
		if strings.HasPrefix(sk, "MEALTYPE#") {
			types = append(types, sk[9:])
		}
	}
	return types, nil
}

// GetAllParticipationForDate returns explicit participation counts per meal type for a given date.
// Result is a map of mealType → map of participation status ("YES"/"NO") → count.
func GetAllParticipationForDate(date string) (map[string]map[string]int, error) {
	db := database.GetDBClient()
	table := database.GetTableName()

	out, err := db.Query(&dynamodb.QueryInput{
		TableName:              aws.String(table),
		KeyConditionExpression: aws.String("PK = :pk AND begins_with(SK, :prefix)"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":pk":     {S: aws.String("DAY#" + date)},
			":prefix": {S: aws.String("TEAM#")},
		},
	})
	if err != nil {
		return nil, err
	}

	// SK format: TEAM#<teamId>#MEAL#<mealType>#USER#<userId>
	result := map[string]map[string]int{}
	for _, item := range out.Items {
		sk := *item["SK"].S
		mealIdx := strings.Index(sk, "#MEAL#")
		userIdx := strings.Index(sk, "#USER#")
		if mealIdx < 0 || userIdx <= mealIdx {
			continue
		}
		mealType := sk[mealIdx+6 : userIdx]
		participation := ""
		if v, ok := item["participation"]; ok && v.S != nil {
			participation = *v.S
		}
		if result[mealType] == nil {
			result[mealType] = map[string]int{}
		}
		result[mealType][participation]++
	}
	return result, nil
}

// CountAllUsers returns the total number of registered users in the system.
func CountAllUsers() (int, error) {
	db := database.GetDBClient()
	table := database.GetTableName()

	var total int64
	err := db.ScanPages(&dynamodb.ScanInput{
		TableName:        aws.String(table),
		FilterExpression: aws.String("SK = :meta AND begins_with(PK, :user)"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":meta": {S: aws.String("META")},
			":user": {S: aws.String("USER#")},
		},
		Select: aws.String(dynamodb.SelectCount),
	}, func(page *dynamodb.ScanOutput, _ bool) bool {
		total += aws.Int64Value(page.Count)
		return true
	})
	return int(total), err
}

// SetMealTypeForDate adds a meal type configuration for a specific date.
func SetMealTypeForDate(date, mealType, adminID string) error {
	db := database.GetDBClient()
	table := database.GetTableName()

	_, err := db.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String(table),
		Item: map[string]*dynamodb.AttributeValue{
			"PK":      {S: aws.String("DAY#" + date)},
			"SK":      {S: aws.String("MEALTYPE#" + mealType)},
			"addedBy": {S: aws.String(adminID)},
		},
	})
	return err
}

// GetWorkLocation returns the stored work location for a user on a given date.
// Returns "OFFICE" if no record exists (opted in by default).
func GetWorkLocation(userID, date string) (string, error) {
	db := database.GetDBClient()
	table := database.GetTableName()

	result, err := db.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(table),
		Key: map[string]*dynamodb.AttributeValue{
			"PK": {S: aws.String("DAY#" + date)},
			"SK": {S: aws.String("LOCATION#USER#" + userID)},
		},
	})
	if err != nil {
		return "", err
	}
	if result.Item == nil {
		return "OFFICE", nil
	}
	if v, ok := result.Item["location"]; ok && v.S != nil {
		return *v.S, nil
	}
	return "OFFICE", nil
}

// GetWorkLocationCountsForDate returns how many users have each work location explicitly set for a date.
// Result is a map of location ("OFFICE"/"WFH") → count of explicit records.
func GetWorkLocationCountsForDate(date string) (map[string]int, error) {
	db := database.GetDBClient()
	table := database.GetTableName()

	out, err := db.Query(&dynamodb.QueryInput{
		TableName:              aws.String(table),
		KeyConditionExpression: aws.String("PK = :pk AND begins_with(SK, :prefix)"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":pk":     {S: aws.String("DAY#" + date)},
			":prefix": {S: aws.String("LOCATION#USER#")},
		},
	})
	if err != nil {
		return nil, err
	}

	counts := map[string]int{}
	for _, item := range out.Items {
		if v, ok := item["location"]; ok && v.S != nil {
			counts[*v.S]++
		}
	}
	return counts, nil
}

// SetWorkLocation stores a user's work location for a given date.
func SetWorkLocation(userID, date, location string) error {
	db := database.GetDBClient()
	table := database.GetTableName()

	_, err := db.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String(table),
		Item: map[string]*dynamodb.AttributeValue{
			"PK":       {S: aws.String("DAY#" + date)},
			"SK":       {S: aws.String("LOCATION#USER#" + userID)},
			"location": {S: aws.String(location)},
		},
	})
	return err
}
