package controllers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ishu17077/code_runner_backend/worker-node/constants"
	"github.com/ishu17077/code_runner_backend/worker-node/helpers"
	"github.com/ishu17077/code_runner_backend/worker-node/middlewares"
)

func AdminLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		bodyAsBytes, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body"})
			return
		}
		var jsonMap map[string]any
		err = json.Unmarshal(bodyAsBytes, &jsonMap)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON Body"})
			return
		}

		adminId, ok := jsonMap["id"]
		if !ok || adminId == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Admin ID not provided"})
			return
		}
		adminPass, ok := jsonMap["password"]
		if !ok || adminPass == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No password is provided"})
			return
		}
		//! Hardcoded env login id and pass
		if constants.AdminLoginId == adminId && constants.AdminLoginPass == adminPass {
			token, refTok, err := helpers.GenerateTokens(adminId.(string))
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Generation of tokens unsuccessful"})
				return
			}
			middlewares.SetCookie(c, "token", token)
			middlewares.SetCookie(c, "refresh_token", refTok)
			c.JSON(http.StatusOK, gin.H{"error": nil, "token": token, "refresh_token": refTok})
			return
		}

		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid id or password"})

	}
}

func AdminLogout() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}
