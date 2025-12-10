package CompositionRoot

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/ishu17077/code_runner_backend/routes"
)

var router *gin.Engine
var server *http.Server

func init() {
	port := os.Getenv("WORKER_PORT")
	if port == "" {
		port = "8060"
	}
	router = gin.New()
	router.Use(gin.Logger())

	server = &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}
}

func Start() error {
	routesDefine(router)
	return server.ListenAndServe()
}

func Stop() error {
	return server.Close()
}

func routesDefine(router *gin.Engine) {
	submissionRoutes := router.Group("/submission")
	adminRoutes := router.Group("/admin")

	routes.SubmissionRoutes(submissionRoutes)
	routes.AdminRoutes(adminRoutes)
}
