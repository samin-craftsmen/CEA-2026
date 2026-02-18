package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/samin-craftsmen/gin-project/handlers"
)

func RegisterMealRoutes(rg *gin.RouterGroup) {

	rg.GET("/meals/today", handlers.GetTodayMeals)
	rg.GET("/meals/tomorrow", handlers.GetTomorrowMeals)

	//rg.POST("/meals/bulk-opt-in", handlers.BulkOptIn)
	//rg.POST("/meals/bulk-opt-out", handlers.BulkOptOut)

	rg.GET("/meals/headcount/:date", handlers.Headcount)

	rg.POST("/meals/items/update", handlers.UpdateMeals)

	rg.GET("/meals/items/:date", handlers.GetMealItemsByDate)

	rg.POST("/meals/override", handlers.OverrideMealSelection)

	rg.POST("/meals/update", handlers.UpdateMealSelection)

	rg.GET("/me", handlers.GetUser)
}
