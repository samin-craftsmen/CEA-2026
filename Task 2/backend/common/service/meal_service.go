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

	mealTypes, err := GetMealTypesForDate(date)
	if err != nil {
		return nil, err
	}

	result := map[string]string{}
	for _, mt := range mealTypes {
		result[mt] = "YES"
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
	if err := validateMealType(mealType, date); err != nil {
		return err
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

type TeamMealViewResponse struct {
	Date    string                       `json:"date"`
	TeamID  string                       `json:"team_id"`
	Members map[string]map[string]string `json:"members"`
}

// defaultMealTypes lists the meal types available for every date.
var defaultMealTypes = []string{"lunch", "snacks"}

// GetMealTypesForDate returns the effective meal types for a date (defaults + date-specific additions).
func GetMealTypesForDate(date string) ([]string, error) {
	extra, err := repository.GetMealTypesForDate(date)
	if err != nil {
		return nil, err
	}

	seen := make(map[string]bool)
	var result []string
	for _, mt := range defaultMealTypes {
		if !seen[mt] {
			seen[mt] = true
			result = append(result, mt)
		}
	}
	for _, mt := range extra {
		if !seen[mt] {
			seen[mt] = true
			result = append(result, mt)
		}
	}
	return result, nil
}

// validateMealType checks if the given meal type is valid for the specified date.
func validateMealType(mealType, date string) error {
	types, err := GetMealTypesForDate(date)
	if err != nil {
		return err
	}
	for _, mt := range types {
		if mt == mealType {
			return nil
		}
	}
	return &ValidationError{fmt.Sprintf("invalid meal type '%s' for date %s: available types are %s", mealType, date, strings.Join(types, ", "))}
}

// AdminAddMealType adds a meal type for a specific date. Only admins can perform this action.
func AdminAddMealType(adminID, date, mealType string) error {
	role, _, err := repository.GetUserRole(adminID)
	if err != nil {
		return err
	}
	if role != "ADMIN" {
		return &ValidationError{"access denied: only admins can perform this action"}
	}

	mealType = strings.ToLower(mealType)
	if mealType == "" {
		return &ValidationError{"meal type cannot be empty"}
	}

	return repository.SetMealTypeForDate(date, mealType, adminID)
}

// TeamLeadGetMealView returns meal participation for all members of the team lead's team on a given date.
// Any meal type not explicitly recorded defaults to YES.
func TeamLeadGetMealView(teamLeadID, date string) (*TeamMealViewResponse, error) {
	role, teamID, err := repository.GetUserRole(teamLeadID)
	if err != nil {
		return nil, err
	}
	if role != "TEAM_LEAD" {
		return nil, &ValidationError{"access denied: only team leads can perform this action"}
	}

	teamMembers, err := repository.GetTeamMembers(teamID)
	if err != nil {
		return nil, err
	}

	existing, err := repository.GetTeamMeals(teamID, date)
	if err != nil {
		return nil, err
	}

	mealTypes, err := GetMealTypesForDate(date)
	if err != nil {
		return nil, err
	}

	// Build result: every member gets every meal type, defaulting to YES.
	result := make(map[string]map[string]string, len(teamMembers))
	for _, memberID := range teamMembers {
		meals := map[string]string{}
		for _, mt := range mealTypes {
			meals[mt] = "YES"
		}
		if recorded, ok := existing[memberID]; ok {
			for mt, status := range recorded {
				meals[mt] = status
			}
		}
		result[memberID] = meals
	}

	return &TeamMealViewResponse{
		Date:    date,
		TeamID:  teamID,
		Members: result,
	}, nil
}

// TeamLeadSetMealStatus updates a team member's meal participation on behalf of a team lead.
// It enforces team-lead-only access, team boundary restrictions, and a 9pm cutoff.
func TeamLeadSetMealStatus(teamLeadID, targetUserID, date, mealType, status string) error {
	role, teamID, err := repository.GetUserRole(teamLeadID)
	if err != nil {
		return err
	}
	if role != "TEAM_LEAD" {
		return &ValidationError{"access denied: only team leads can perform this action"}
	}

	isMember, err := repository.VerifyTeamMembership(teamID, targetUserID)
	if err != nil {
		return err
	}
	if !isMember {
		return &ValidationError{"access denied: target user does not belong to your team"}
	}

	mealType = strings.ToLower(mealType)
	if err := validateMealType(mealType, date); err != nil {
		return err
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

	return repository.SetMealParticipation(targetUserID, date, mealType, status, teamID)
}

type AdminMealViewResponse struct {
	Date   string            `json:"date"`
	UserID string            `json:"user_id"`
	Meals  map[string]string `json:"meals"`
}

// AdminGetMealView returns meal participation for a specific employee on a given date.
// Only admins can use this. Missing records default to YES (opted-in by default).
func AdminGetMealView(adminID, targetUserID, date string) (*AdminMealViewResponse, error) {
	role, _, err := repository.GetUserRole(adminID)
	if err != nil {
		return nil, err
	}
	if role != "ADMIN" {
		return nil, &ValidationError{"access denied: only admins can perform this action"}
	}

	meals, err := repository.GetUserMeals(targetUserID, date)
	if err != nil {
		return nil, err
	}

	mealTypes, err := GetMealTypesForDate(date)
	if err != nil {
		return nil, err
	}

	result := map[string]string{}
	for _, mt := range mealTypes {
		result[mt] = "YES"
	}
	for mealType, participation := range meals {
		result[mealType] = participation
	}

	return &AdminMealViewResponse{
		Date:   date,
		UserID: targetUserID,
		Meals:  result,
	}, nil
}

// AdminSetMealStatus updates any employee's meal participation on behalf of an admin.
// It enforces admin-only access and a 9pm cutoff for the following day's meals.
func AdminSetMealStatus(adminID, targetUserID, date, mealType, status string) error {
	role, _, err := repository.GetUserRole(adminID)
	if err != nil {
		return err
	}
	if role != "ADMIN" {
		return &ValidationError{"access denied: only admins can perform this action"}
	}

	mealType = strings.ToLower(mealType)
	if err := validateMealType(mealType, date); err != nil {
		return err
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

	_, targetTeamID, err := repository.GetUserRole(targetUserID)
	if err != nil {
		return err
	}

	return repository.SetMealParticipation(targetUserID, date, mealType, status, targetTeamID)
}

type HeadcountEntry struct {
	Yes int `json:"yes"`
	No  int `json:"no"`
}

type WorkLocationSummary struct {
	Office int `json:"office"`
	WFH    int `json:"wfh"`
}

type HeadcountSummaryResponse struct {
	Date         string                    `json:"date"`
	WorkLocation WorkLocationSummary       `json:"work_location"`
	Summary      map[string]HeadcountEntry `json:"summary"`
}

// AdminGetHeadcountSummary returns per-meal-type headcount for a given date.
// Only admins may call this.
// Because meals are opted-in by default, YES = total registered users - explicit NO count.
func AdminGetHeadcountSummary(adminID, date string) (*HeadcountSummaryResponse, error) {
	role, _, err := repository.GetUserRole(adminID)
	if err != nil {
		return nil, err
	}
	if role != "ADMIN" {
		return nil, &ValidationError{"access denied: only admins can perform this action"}
	}

	totalUsers, err := repository.CountAllUsers()
	if err != nil {
		return nil, err
	}

	counts, err := repository.GetAllParticipationForDate(date)
	if err != nil {
		return nil, err
	}

	mealTypes, err := GetMealTypesForDate(date)
	if err != nil {
		return nil, err
	}

	summary := make(map[string]HeadcountEntry, len(mealTypes))
	for _, mt := range mealTypes {
		noCount := 0
		if c, ok := counts[mt]; ok {
			noCount = c["NO"]
		}
		summary[mt] = HeadcountEntry{
			Yes: totalUsers - noCount,
			No:  noCount,
		}
	}

	locationCounts, err := repository.GetWorkLocationCountsForDate(date)
	if err != nil {
		return nil, err
	}
	// Users with no location record are implicitly OFFICE.
	explicitWFH := locationCounts["WFH"]
	explicitOffice := locationCounts["OFFICE"]
	workLocation := WorkLocationSummary{
		Office: totalUsers - explicitWFH,
		WFH:    explicitWFH,
	}
	_ = explicitOffice // counted implicitly above

	return &HeadcountSummaryResponse{
		Date:         date,
		WorkLocation: workLocation,
		Summary:      summary,
	}, nil
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

// WorkLocationResponse is returned when viewing a user's work location for a date.
type WorkLocationResponse struct {
	Date     string `json:"date"`
	UserID   string `json:"user_id"`
	Location string `json:"location"`
}

// GetWorkLocationForDate returns a user's work location for a given date.
// Defaults to OFFICE if no record exists.
func GetWorkLocationForDate(discordID, date string) (*WorkLocationResponse, error) {
	if err := repository.EnsureUserExists(discordID); err != nil {
		return nil, err
	}

	location, err := repository.GetWorkLocation(discordID, date)
	if err != nil {
		return nil, err
	}

	return &WorkLocationResponse{
		Date:     date,
		UserID:   discordID,
		Location: location,
	}, nil
}

// SetWorkLocationForDate updates a user's work location for a given date.
// Switching to WFH opts out all meals; switching to OFFICE opts all meals back in.
func SetWorkLocationForDate(discordID, date, location string) error {
	if err := repository.EnsureUserExists(discordID); err != nil {
		return err
	}

	_, teamID, err := repository.GetUserRole(discordID)
	if err != nil {
		return err
	}

	location = strings.ToUpper(location)
	if location != "OFFICE" && location != "WFH" {
		return &ValidationError{fmt.Sprintf("invalid location '%s': must be 'OFFICE' or 'WFH'", location)}
	}

	locked, err := isPastCutoff(date)
	if err != nil {
		return &ValidationError{err.Error()}
	}
	if locked {
		return &ValidationError{"the cutoff time (9pm) has passed — work location can no longer be updated for this date"}
	}

	mealTypes, err := GetMealTypesForDate(date)
	if err != nil {
		return err
	}

	status := "YES"
	if location == "WFH" {
		status = "NO"
	}

	for _, mt := range mealTypes {
		if err := repository.SetMealParticipation(discordID, date, mt, status, teamID); err != nil {
			return err
		}
	}

	return repository.SetWorkLocation(discordID, date, location)
}

// Day status type constants.
const (
	DayStatusGovernmentHoliday = "GOVERNMENT_HOLIDAY"
	DayStatusOfficeClosed      = "OFFICE_CLOSED"
	DayStatusSpecialEvent      = "SPECIAL_EVENT"
)

// DayStatusResponse is returned when viewing the status of a specific day.
type DayStatusResponse struct {
	Date  string `json:"date"`
	Type  string `json:"type"`
	Note  string `json:"note,omitempty"`
	SetBy string `json:"set_by,omitempty"`
}

// AdminSetDayStatus sets the administrative status for a specific day. Only admins may call this.
// GOVERNMENT_HOLIDAY and SPECIAL_EVENT only update the day status marker.
// OFFICE_CLOSED additionally opts out all registered users from every meal on that date.
// A non-empty note is required for SPECIAL_EVENT.
func AdminSetDayStatus(adminID, date, statusType, note string) error {
	role, _, err := repository.GetUserRole(adminID)
	if err != nil {
		return err
	}
	if role != "ADMIN" {
		return &ValidationError{"access denied: only admins can perform this action"}
	}

	switch statusType {
	case DayStatusGovernmentHoliday, DayStatusOfficeClosed, DayStatusSpecialEvent:
		// valid
	default:
		return &ValidationError{fmt.Sprintf("invalid status type '%s': must be GOVERNMENT_HOLIDAY, OFFICE_CLOSED, or SPECIAL_EVENT", statusType)}
	}

	if statusType == DayStatusSpecialEvent && strings.TrimSpace(note) == "" {
		return &ValidationError{"a note is required for SPECIAL_EVENT day status"}
	}

	if statusType == DayStatusOfficeClosed {
		if err := optOutAllMealsForDate(date); err != nil {
			return err
		}
	}

	return repository.SetDayStatus(date, statusType, note, adminID)
}

// optOutAllMealsForDate sets all meal participation records to NO for every registered user on the given date.
func optOutAllMealsForDate(date string) error {
	allUsers, err := repository.GetAllUserIDs()
	if err != nil {
		return err
	}

	mealTypes, err := GetMealTypesForDate(date)
	if err != nil {
		return err
	}

	for _, userID := range allUsers {
		_, teamID, err := repository.GetUserRole(userID)
		if err != nil {
			continue // skip users without a valid role/team record
		}
		for _, mt := range mealTypes {
			if err := repository.SetMealParticipation(userID, date, mt, "NO", teamID); err != nil {
				return err
			}
		}
	}
	return nil
}

// GetHeadcountSummaryForDate returns the headcount summary for a date without role enforcement.
// Intended for internal/scheduled use only (e.g. the nightly EventBridge notification).
func GetHeadcountSummaryForDate(date string) (*HeadcountSummaryResponse, error) {
	totalUsers, err := repository.CountAllUsers()
	if err != nil {
		return nil, err
	}

	counts, err := repository.GetAllParticipationForDate(date)
	if err != nil {
		return nil, err
	}

	mealTypes, err := GetMealTypesForDate(date)
	if err != nil {
		return nil, err
	}

	summary := make(map[string]HeadcountEntry, len(mealTypes))
	for _, mt := range mealTypes {
		noCount := 0
		if c, ok := counts[mt]; ok {
			noCount = c["NO"]
		}
		summary[mt] = HeadcountEntry{
			Yes: totalUsers - noCount,
			No:  noCount,
		}
	}

	locationCounts, err := repository.GetWorkLocationCountsForDate(date)
	if err != nil {
		return nil, err
	}
	explicitWFH := locationCounts["WFH"]

	return &HeadcountSummaryResponse{
		Date:         date,
		WorkLocation: WorkLocationSummary{Office: totalUsers - explicitWFH, WFH: explicitWFH},
		Summary:      summary,
	}, nil
}

// GetDayStatusForDate returns the administrative status for a given date.
// Returns Type="NORMAL" when no status has been set.
func GetDayStatusForDate(date string) (*DayStatusResponse, error) {
	ds, err := repository.GetDayStatus(date)
	if err != nil {
		return nil, err
	}
	if ds == nil {
		return &DayStatusResponse{Date: date, Type: "NORMAL"}, nil
	}
	return &DayStatusResponse{
		Date:  date,
		Type:  ds.Type,
		Note:  ds.Note,
		SetBy: ds.SetBy,
	}, nil
}
