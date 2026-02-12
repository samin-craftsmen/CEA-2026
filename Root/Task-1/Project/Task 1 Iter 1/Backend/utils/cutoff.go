package utils

import "time"

func IsBeforeCutoff() bool {
	now := time.Now()
	today := now.Format("2006-01-02")

	loc := now.Location()
	cutoffTime, _ := time.ParseInLocation(
		"2006-01-02 15:04",
		today+" 21:00",
		loc,
	)

	return now.Before(cutoffTime)
}
