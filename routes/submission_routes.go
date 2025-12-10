package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/ishu17077/code_runner_backend/controllers"
)

func SubmissionRoutes(incomingRoutes *gin.RouterGroup) {
	incomingRoutes.POST("/test/public", controllers.PublicTestSubmission())
	incomingRoutes.POST("/test/private", controllers.PrivateTestSubmission())
}
