package main

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/samin-craftsmen/gin-project/middleware"
	"github.com/samin-craftsmen/gin-project/routes"
)

func main() {
	r := gin.Default()

	// CORS configuration
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Public routes
	routes.RegisterAuthRoutes(r)

	// Protected routes
	authorized := r.Group("/")
	authorized.Use(middleware.AuthMiddleware())

	routes.RegisterMealRoutes(authorized)
	routes.RegisterAdminRoutes(authorized)
	routes.RegisterTeamRoutes(authorized)
	routes.RegisterMeRoutes(authorized)

	r.Run(":8080")
}
