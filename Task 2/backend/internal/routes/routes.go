package routes

import (
	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all API routes
func SetupRoutes(router *gin.Engine) {
	// Health check endpoint
    router.GET("/health", handlers.HealthCheck)
}
