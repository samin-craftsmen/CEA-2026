package service

import (
	"github.com/samin-craftsmen/meal-headcount-planner-backend/common/repository"
)

type MealViewResponse struct {
	Date   string            `json:"date"`
	Meals  map[string]string `json:"meals"`
	UserID string            `json:"user_id"`
}

func GetMealView(discordID string, date string) (*MealViewResponse, error) {
	if err := repository.EnsureUserExists(discordID); err != nil {
		return nil, err
	}

	meals, err := repository.GetUserMeals(discordID, date)
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
		Date:   date,
		Meals:  result,
		UserID: discordID,
	}, nil
}
