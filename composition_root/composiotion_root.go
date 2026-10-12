package CompositionRoot

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/ishu17077/code_runner_backend/middlewares"
	"github.com/ishu17077/code_runner_backend/routes"
)

var router *gin.Engine
var server *http.Server
var upgrader *websocket.Upgrader

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

	upgrader = &websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

}

func Start() error {
	middlewareUse(router)
	routesDefine(router)

	return server.ListenAndServe()
}

func Stop() error {
	return server.Close()
}

func routesDefine(router *gin.Engine) {
	submissionRoutes := router.Group("/submission")

	adminRoutes := router.Group("/admin")

	routes.SubmissionRoutes(submissionRoutes, upgrader)
	routes.AdminRoutes(adminRoutes)
}

func middlewareUse(router *gin.Engine) {
	router.Use(middlewares.CORSMiddleware())
	router.SetTrustedProxies([]string{"127.0.0.1", "localhost"})
	router.Use(middlewares.MaxAllowedSize(5 * 1024 * 1024)) //? 5 MB
}
