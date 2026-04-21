package commands

import (
	"fmt"
	"strings"

	"github.com/samin-craftsmen/meal-headcount-planner-backend/common/validation"
	"github.com/samin-craftsmen/meal-headcount-planner-backend/lambdas/google-chat/types"
)

func requireExactArgs(parts []string, expected int, usage string) *types.Response {
	if len(parts) != expected {
		return types.ErrorResponse("usage: " + usage)
	}
	return nil
}

func requireMinArgs(parts []string, minimum int, usage string) *types.Response {
	if len(parts) < minimum {
		return types.ErrorResponse("usage: " + usage)
	}
	return nil
}

func validatedCommandUserID(userID string) (string, *types.Response) {
	validated, err := validation.UserID(userID, "user id")
	if err != nil {
		return "", types.ErrorResponse(err.Error())
	}
	return validated, nil
}

func validatedChatDate(value string) (string, *types.Response) {
	validated, err := validation.Date(value)
	if err != nil {
		return "", types.ErrorResponse(err.Error())
	}
	return validated, nil
}

func validatedChatMealType(value string) (string, *types.Response) {
	validated, err := validation.MealType(value)
	if err != nil {
		return "", types.ErrorResponse(err.Error())
	}
	return validated, nil
}

func validatedChatStatus(value string) (string, *types.Response) {
	validated, err := validation.Status(value)
	if err != nil {
		return "", types.ErrorResponse(err.Error())
	}
	return validated, nil
}

func validatedChatLocation(value string) (string, *types.Response) {
	validated, err := validation.Location(value)
	if err != nil {
		return "", types.ErrorResponse(err.Error())
	}
	return validated, nil
}

func validatedChatDayStatus(value string) (string, *types.Response) {
	validated, err := validation.DayStatusType(value)
	if err != nil {
		return "", types.ErrorResponse(err.Error())
	}
	return validated, nil
}

func validatedChatNote(value string, required bool) (string, *types.Response) {
	validated, err := validation.Note(value, required)
	if err != nil {
		return "", types.ErrorResponse(err.Error())
	}
	return validated, nil
}

func validatedChatTargetUserID(value string) (string, *types.Response) {
	normalized := normalizeUserRef(value)
	validated, err := validation.UserID(normalized, "user_id")
	if err != nil {
		return "", types.ErrorResponse(err.Error())
	}
	return validated, nil
}

func normalizedSubcommand(value string) string {
	return strings.ToLower(strings.TrimSpace(value))
}

func invalidStringValue(name string) *types.Response {
	return types.ErrorResponse(fmt.Sprintf("invalid %s value", name))
}
