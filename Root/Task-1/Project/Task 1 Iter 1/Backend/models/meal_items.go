package models

type MealItems struct {
	Date string `json:"date"` // YYYY-MM-DD

	Items map[MealType][]string `json:"items"`
}
