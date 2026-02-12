package utils

import (
	"encoding/json"
	"os"

	"github.com/yourname/gin-project/models"
)

const mealItemsFile = "meal_items.json"

func LoadMealItems() ([]models.MealItems, error) {
	file, err := os.ReadFile(mealItemsFile)
	if err != nil {
		return []models.MealItems{}, nil
	}

	var data []models.MealItems
	json.Unmarshal(file, &data)

	return data, nil
}

func SaveMealItems(data []models.MealItems) error {
	bytes, _ := json.MarshalIndent(data, "", "  ")
	return os.WriteFile(mealItemsFile, bytes, 0644)
}
