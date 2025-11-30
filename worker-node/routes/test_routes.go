package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/ishu17077/code_runner_backend/worker-node/controllers"
	"github.com/ishu17077/code_runner_backend/worker-node/middlewares"
)

func TestRoutes(incomingRoutes *gin.RouterGroup) {
	incomingRoutes.Use(middlewares.Authenticate())
	incomingRoutes.GET("/", controllers.GetTests())
	incomingRoutes.GET("/:test_id", controllers.GetTest())
	incomingRoutes.POST("/", controllers.CreateTest())
	incomingRoutes.PATCH("/:test_id", controllers.UpdateTest())
	//TODO: Post and patch routes
}
