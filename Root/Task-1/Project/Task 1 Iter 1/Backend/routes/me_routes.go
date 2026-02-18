package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/samin-craftsmen/gin-project/handlers"
)

func RegisterMeRoutes(rg *gin.RouterGroup) {
	rg.GET("/me", handlers.GetUser)
	rg.POST("/me/work-location", handlers.SetWorkLocation)
	rg.GET("me/work-location", handlers.GetWorkLocation)
}
