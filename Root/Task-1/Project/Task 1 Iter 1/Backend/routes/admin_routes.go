package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/samin-craftsmen/gin-project/handlers"
)

func RegisterAdminRoutes(rg *gin.RouterGroup) {

	rg.POST("/admin/meals/opt-in/:date", handlers.AdminOptIn)
	rg.POST("/admin/meals/opt-out/:date", handlers.AdminOptOut)

	rg.GET("/admin/teams/meals/:date", handlers.GetTeamMealCountsByDate)

	rg.POST("admin/day-controls", handlers.SetSpecialDay)
	rg.GET("admin/day-controls/:date", handlers.GetDayStatus)
}
