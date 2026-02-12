package models

type MealType string

const (
	Lunch          MealType = "Lunch"
	Snacks         MealType = "Snacks"
	Iftar          MealType = "Iftar"
	EventDinner    MealType = "Event Dinner"
	OptionalDinner MealType = "Optional Dinner"
)

type Participation struct {
	Username string     `json:"username"`
	Date     string     `json:"date"`  // YYYY-MM-DD
	Meals    []MealType `json:"meals"` // meals opted IN
}
