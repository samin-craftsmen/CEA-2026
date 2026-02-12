package utils

import (
	"encoding/json"
	"os"

	"github.com/yourname/gin-project/models"
)

const participationFile = "participation.json"

func LoadParticipation() ([]models.Participation, error) {
	file, err := os.ReadFile(participationFile)
	if err != nil {
		return nil, err
	}

	var data []models.Participation
	err = json.Unmarshal(file, &data)
	return data, err
}

func SaveParticipation(data []models.Participation) error {
	bytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(participationFile, bytes, 0644)
}
