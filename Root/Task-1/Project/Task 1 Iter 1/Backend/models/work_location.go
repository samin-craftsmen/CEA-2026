package models

type WorkLocation struct {
	Username string `json:"username"`
	Date     string `json:"date"`
	Location string `json:"location"` // "Office" or "WFH"
}
