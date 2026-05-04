package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/ishu17077/code_runner_backend/controllers"
)

func SubmissionRoutes(incomingRoutes *gin.RouterGroup, upgrader *websocket.Upgrader) {

	incomingRoutes.POST("/test/public", controllers.PublicTestSubmission())
	incomingRoutes.POST("/test/private", controllers.PrivateTestSubmission())
	incomingRoutes.GET("/run", controllers.RunCode(upgrader))
}
