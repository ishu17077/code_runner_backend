package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/ishu17077/code_runner_backend/worker-node/controllers"
)

func TestRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("/initial-test", controllers.InitialTest())
}