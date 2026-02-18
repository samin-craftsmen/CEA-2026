package models

type DayControl struct {
	Date string `json:"date"`

	Type string `json:"type"`
	// allowed values:
	// "office_closed"
	// "government_holiday"
	// "special_celebration"

	Note *string `json:"note,omitempty"`
}
