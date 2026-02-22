package utils

import (
	"encoding/json"
	"os"

	"github.com/samin-craftsmen/gin-project/models"
)

func LoadWorkLocations() ([]models.WorkLocation, error) {

	file, err := os.ReadFile("work_locations.json")
	if err != nil {
		return nil, err
	}

	var locations []models.WorkLocation

	if len(file) == 0 {
		return locations, nil
	}

	err = json.Unmarshal(file, &locations)
	if err != nil {
		return nil, err
	}

	return locations, nil
}

func SaveWorkLocations(locations []models.WorkLocation) error {

	data, err := json.MarshalIndent(locations, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile("work_locations.json", data, 0644)
}
