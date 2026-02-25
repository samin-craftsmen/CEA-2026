package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/samin-craftsmen/gin-project/handlers"
)

func RegisterTeamRoutes(rg *gin.RouterGroup) {

	rg.GET("/teams/meals/today", handlers.GetTodayTeamMeals)
	rg.POST("/teams/meals/optout", handlers.TeamBulkOptOut)
	rg.POST("/teams/meals/optin", handlers.TeamBulkOptIn)
	rg.POST("/teams/work-location/update", handlers.UpdateTeamMemberWorkLocation)
	rg.POST("/teams/work-location", handlers.GetTeamMemberWorkLocation)
}
