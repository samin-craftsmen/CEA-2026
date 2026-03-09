package service

import (
	"time"

	"github.com/samin-craftsmen/meal-headcount-planner-backend/common/repository"
)

type MealViewResponse struct {
	Date   string            `json:"date"`
	Meals  map[string]string `json:"meals"`
	UserID string            `json:"user_id"`
}

func GetMealView(discordID string) (*MealViewResponse, error) {
	if err := repository.EnsureUserExists(discordID); err != nil {
		return nil, err
	}

	tomorrow := time.Now().Add(24 * time.Hour).Format("2006-01-02")

	meals, err := repository.GetUserMeals(discordID, tomorrow)
	if err != nil {
		return nil, err
	}

	result := map[string]string{
		"lunch":  "YES",
		"snacks": "YES",
	}
	for mealType, participation := range meals {
		result[mealType] = participation
	}

	return &MealViewResponse{
		Date:   tomorrow,
		Meals:  result,
		UserID: discordID,
	}, nil
}
