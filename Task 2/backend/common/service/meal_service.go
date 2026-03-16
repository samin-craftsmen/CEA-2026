package service

import (
	"fmt"
	"strings"
	"time"

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

// ValidationError represents a user-facing validation failure (e.g. bad input, cutoff passed).
type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string { return e.Message }

// SetMealStatus updates a user's participation for a specific meal type on a given date.
// It enforces employee-only access and a 9pm cutoff for the following day's meals.
func SetMealStatus(discordID, date, mealType, status string) error {
	if err := repository.EnsureUserExists(discordID); err != nil {
		return err
	}

	_, teamID, err := repository.GetUserRole(discordID)
	if err != nil {
		return err
	}

	mealType = strings.ToLower(mealType)
	if mealType != "lunch" && mealType != "snacks" {
		return &ValidationError{fmt.Sprintf("invalid meal type '%s': must be 'lunch' or 'snacks'", mealType)}
	}

	status = strings.ToUpper(status)
	if status != "YES" && status != "NO" {
		return &ValidationError{fmt.Sprintf("invalid status '%s': must be 'YES' or 'NO'", status)}
	}

	locked, err := isPastCutoff(date)
	if err != nil {
		return &ValidationError{err.Error()}
	}
	if locked {
		return &ValidationError{"the cutoff time (9pm) has passed — meal status can no longer be updated for this date"}
	}

	return repository.SetMealParticipation(discordID, date, mealType, status, teamID)
}

// isPastCutoff reports whether the cutoff for updating the given date has passed.
// Rules (all times in IST / UTC+5:30):
//   - Target date is today or in the past → always locked.
//   - Target date is tomorrow → locked once it is past 9pm today.
//   - Target date is two or more days away → not locked.
func isPastCutoff(date string) (bool, error) {
	targetDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		return false, fmt.Errorf("invalid date format, expected YYYY-MM-DD")
	}

	ist := time.FixedZone("IST", 5*60*60+30*60)
	now := time.Now().In(ist)

	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, ist)
	tomorrow := today.AddDate(0, 0, 1)
	target := time.Date(targetDate.Year(), targetDate.Month(), targetDate.Day(), 0, 0, 0, 0, ist)

	// Today or in the past
	if !target.After(today) {
		return true, nil
	}

	// Tomorrow: check 9pm cutoff
	if target.Equal(tomorrow) {
		cutoff := time.Date(now.Year(), now.Month(), now.Day(), 21, 0, 0, 0, ist)
		return now.After(cutoff), nil
	}

	// Day after tomorrow or further: not locked
	return false, nil
}
