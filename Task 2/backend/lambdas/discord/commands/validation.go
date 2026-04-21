package commands

import (
	"fmt"
	"strings"

	"github.com/samin-craftsmen/meal-headcount-planner-backend/common/validation"
	"github.com/samin-craftsmen/meal-headcount-planner-backend/lambdas/discord/types"
)

func commandUserID(userID string) (string, *types.InteractionResponse) {
	validated, err := validation.UserID(userID, "user id")
	if err != nil {
		return "", types.ErrorResponse(err.Error())
	}
	return validated, nil
}

func stringOption(options []types.CommandOption, name string) (string, error) {
	for _, opt := range options {
		if opt.Name != name {
			continue
		}
		if opt.Value == nil {
			return "", nil
		}
		value, ok := opt.Value.(string)
		if !ok {
			return "", fmt.Errorf("invalid %s option: expected a string", name)
		}
		return value, nil
	}
	return "", nil
}

func validatedDate(value string) (string, *types.InteractionResponse) {
	validated, err := validation.Date(value)
	if err != nil {
		return "", types.ErrorResponse(err.Error())
	}
	return validated, nil
}

func validatedMealType(value string) (string, *types.InteractionResponse) {
	validated, err := validation.MealType(value)
	if err != nil {
		return "", types.ErrorResponse(err.Error())
	}
	return validated, nil
}

func validatedStatus(value string) (string, *types.InteractionResponse) {
	validated, err := validation.Status(value)
	if err != nil {
		return "", types.ErrorResponse(err.Error())
	}
	return validated, nil
}

func validatedLocation(value string) (string, *types.InteractionResponse) {
	validated, err := validation.Location(value)
	if err != nil {
		return "", types.ErrorResponse(err.Error())
	}
	return validated, nil
}

func validatedDayStatus(value string) (string, *types.InteractionResponse) {
	validated, err := validation.DayStatusType(value)
	if err != nil {
		return "", types.ErrorResponse(err.Error())
	}
	return validated, nil
}

func validatedNote(value string, required bool) (string, *types.InteractionResponse) {
	validated, err := validation.Note(value, required)
	if err != nil {
		return "", types.ErrorResponse(err.Error())
	}
	return validated, nil
}

func validatedTargetUserID(value string) (string, *types.InteractionResponse) {
	validated, err := validation.UserID(strings.TrimSpace(value), "employee")
	if err != nil {
		return "", types.ErrorResponse(err.Error())
	}
	return validated, nil
}
