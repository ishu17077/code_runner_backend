package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/ishu17077/code_runner_backend/middlewares"
	"github.com/ishu17077/code_runner_backend/routes"
)

func main() {
	port := os.Getenv("RUNNER_PORT")
	if port == "" {
		port = "8080"
	}
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(middlewares.CORSMiddleware())
	routes.ProgramRoutes(router)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			fmt.Printf("Error: %s\n", err.Error())
			panic(err)
		}
	}()
	//TODO: Impl Shut down gracefully
}
