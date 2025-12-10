package controllers

import (
	"encoding/base64"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/ishu17077/code_runner_backend/helpers/k8s"
	"github.com/ishu17077/code_runner_backend/models"
)

// ! Simple odd even Test

//TODO: Implement sync.Mutex to handle process flow

func PublicTestSubmission() gin.HandlerFunc {
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
		if len(submission.Tests) > 3 {
			c.JSON(http.StatusNotAcceptable, gin.H{"error": "Cannot have more than 3 tests in public submission"})
			return
		}
		codeBytes, err := base64.StdEncoding.DecodeString(submission.Code)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "The provided code is not properly base64 encoded."})
			return
		}

		submission.Code = string(codeBytes)

		res, err := k8s.K8sMgr.RunOnPod(submission)

		if err != nil {
			res.Error = err.Error()
			c.JSON(http.StatusNotAcceptable, res)
			return
		}
		c.JSON(http.StatusAccepted, res)
		// allOkay, execResults, err := helpers.AnalyzeSubmission(submission, submission.Tests)
		// if err != nil {
		// 	c.JSON(http.StatusNotAcceptable, gin.H{"All tests passed": allOkay, "Execution Result": execResults, "Error": err.Error()})
		// 	return
		// }
		// c.JSON(http.StatusOK, gin.H{"All tests passed": allOkay, "Execution Result": execResults})
	}
}

var lock sync.Mutex

func PrivateTestSubmission() gin.HandlerFunc {
	return func(c *gin.Context) {
		lock.Lock()
		defer lock.Unlock()
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

		res, err := k8s.K8sMgr.RunOnPod(submission)

		if err != nil {
			res.Error = err.Error()
			c.JSON(http.StatusNotAcceptable, res)
			return
		}
		c.JSON(http.StatusAccepted, res)
		// allOkay, execResults, err := helpers.AnalyzeSubmission(submission, submission.Tests)
		// if err != nil {
		// 	c.JSON(http.StatusNotAcceptable, gin.H{"All tests passed": allOkay, "Execution Result": execResults, "Error": err.Error()})
		// 	return
		// }
		// c.JSON(http.StatusOK, gin.H{"All tests passed": allOkay, "Execution Result": execResults})
	}
}
