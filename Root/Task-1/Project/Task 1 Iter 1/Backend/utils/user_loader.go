package utils

import (
	"encoding/json"
	"os"

	"github.com/samin-craftsmen/gin-project/models"
)

func LoadUsers() ([]models.User, error) {
	file, err := os.ReadFile("users.json")
	if err != nil {
		return nil, err
	}

	var users []models.User
	err = json.Unmarshal(file, &users)
	return users, err
}
