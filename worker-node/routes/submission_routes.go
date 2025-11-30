package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/ishu17077/code_runner_backend/worker-node/controllers"
)

func SubmissionRoutes(incomingRoutes *gin.RouterGroup) {
	incomingRoutes.POST("/test", controllers.InitialTest())
}
