package utils

import (
	"encoding/json"
	"os"

	"github.com/samin-craftsmen/gin-project/models"
)

var dayControlFile = "day_controls.json"

func LoadDayControls() ([]models.DayControl, error) {

	file, err := os.ReadFile(dayControlFile)
	if err != nil {
		return nil, err
	}

	var data []models.DayControl
	json.Unmarshal(file, &data)

	return data, nil
}

func SaveDayControls(data []models.DayControl) error {

	bytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(dayControlFile, bytes, 0644)
}
