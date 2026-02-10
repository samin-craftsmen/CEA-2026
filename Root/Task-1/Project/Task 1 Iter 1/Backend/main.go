package main

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/yourname/gin-project/middleware"
	"github.com/yourname/gin-project/models"
	"github.com/yourname/gin-project/utils"
)

func main() {
	r := gin.Default()

	// CORS config
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// 🔓 PUBLIC LOGIN ROUTE (returns JWT)
	r.POST("/login", func(c *gin.Context) {
		var loginUser models.User
		if err := c.BindJSON(&loginUser); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		users, _ := utils.LoadUsers()

		for _, user := range users {
			if user.Username == loginUser.Username && user.Password == loginUser.Password {
				token, err := utils.GenerateToken(user.Username)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Token error"})
					return
				}

				c.JSON(http.StatusOK, gin.H{"token": token})
				return
			}
		}

		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
	})

	// 🔒 PROTECTED ROUTES (require JWT)
	authorized := r.Group("/")
	authorized.Use(middleware.AuthMiddleware())
	{
		authorized.GET("/me", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "You are authenticated"})
		})

		// Future meal planner routes go here
		// authorized.GET("/meals", ...)
		// authorized.POST("/meals", ...)
	}

	r.Run(":8080")
}
