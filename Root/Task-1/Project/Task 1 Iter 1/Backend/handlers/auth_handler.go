package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/samin-craftsmen/gin-project/models"
	"github.com/samin-craftsmen/gin-project/utils"
)

func Login(c *gin.Context) {
	var loginUser models.User
	if err := c.BindJSON(&loginUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	users, _ := utils.LoadUsers()

	for _, user := range users {
		if user.Username == loginUser.Username && user.Password == loginUser.Password {
			token, err := utils.GenerateToken(user.Username, user.Role, user.Team)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Token error"})
				return
			}

			c.JSON(http.StatusOK, gin.H{"token": token})
			return
		}
	}

	c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
}
