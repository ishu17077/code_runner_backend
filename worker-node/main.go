package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/ishu17077/code_runner_backend/worker-node/middlewares"
	"github.com/ishu17077/code_runner_backend/worker-node/routes"
)

func main() {
	port := os.Getenv("WORKER_PORT")
	if port == "" {
		port = "8060"
	}

	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())
	router.Use(middlewares.CORSMiddleware())
	routes.TestRoutes(router)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	if err := server.ListenAndServe(); err != nil {
		fmt.Printf("Error starting the server: %s\n", err.Error())
		panic(err)
	}
}
