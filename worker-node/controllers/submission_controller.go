package controllers

import (
	"encoding/base64"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ishu17077/code_runner_backend/worker-node/helpers"
	"github.com/ishu17077/code_runner_backend/worker-node/models"
)

// ! Simple odd even Test

//TODO: Implement sync.Mutex to handle process flow

func InitialTest() gin.HandlerFunc {
	return func(c *gin.Context) {
		var submission models.Submission

		if err := c.ShouldBind(&submission); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request", "msg": err.Error()})
			return
		}

		if len(submission.Tests) == 0 {
			c.JSON(http.StatusNoContent, gin.H{"error": "No tests provided"})
			return
		}
		codeBytes, err := base64.StdEncoding.DecodeString(submission.Code)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "The provided code is not properly base64 encoded."})
			return
		}

		submission.Code = string(codeBytes)
		allOkay, execResults, err := helpers.AnalyzeSubmission(submission, submission.Tests)
		if err != nil {
			c.JSON(http.StatusNotAcceptable, gin.H{"All tests passed": allOkay, "Execution Result": execResults, "Error": err.Error()})
			return
		}
		c.JSON(http.StatusAccepted, gin.H{"All tests passed": allOkay, "Execution Result": execResults})
	}
}

func PrivateTestSubmission() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}
