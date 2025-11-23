package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/ishu17077/code_runner_backend/controllers"
)

func ProgramRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("/submit", controllers.ProgramSubmit())
	incomingRoutes.POST("/validate", controllers.ProgramValidate())
}
