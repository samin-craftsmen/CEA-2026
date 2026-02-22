package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/samin-craftsmen/gin-project/handlers"
)

func RegisterAuthRoutes(r *gin.Engine) {
	r.POST("/login", handlers.Login)
	//r.POST("/logout", handlers.Logout)
	//r.POST("/register", handlers.Register)
}
