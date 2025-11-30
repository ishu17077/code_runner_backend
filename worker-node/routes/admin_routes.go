package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/ishu17077/code_runner_backend/worker-node/controllers"
)

func AdminRoutes(incomingRoutes *gin.RouterGroup) {
	incomingRoutes.POST("/login", controllers.AdminLogin())
	// incomingRoutes.Use(middlewares.Authenticate())
	incomingRoutes.POST("/logout", controllers.AdminLogout())
}
