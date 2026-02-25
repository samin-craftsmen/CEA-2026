package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/samin-craftsmen/gin-project/handlers"
)

func RegisterMeRoutes(rg *gin.RouterGroup) {
	rg.GET("/me", handlers.GetUser)
	rg.POST("/me/work-location", handlers.SetWorkLocation)
	rg.GET("me/work-location", handlers.GetWorkLocation)
	rg.GET("/me/wfh-count", handlers.MeWFHCountHandler)
	// Add this line to your meal routes (or me_routes if that's where user endpoints are)
	rg.GET("/me/meals/today", handlers.GetTodayMeals)
}
