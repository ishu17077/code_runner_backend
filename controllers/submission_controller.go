package controllers

import (
	"encoding/base64"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/ishu17077/code_runner_backend/helpers/k8s"
	"github.com/ishu17077/code_runner_backend/models"
	currentstatus "github.com/ishu17077/code_runner_backend/models/enums/current_status"
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
		if len(submission.Tests) > 5 {
			c.JSON(http.StatusNotAcceptable, gin.H{"error": "Cannot have more than 3 tests in public submission"})
			return
		}
		codeBytes, err := base64.StdEncoding.DecodeString(submission.Code)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "The provided code is not properly base64 encoded."})
			return
		}
		if err := models.ValidateSubmission(&submission); err != nil {
			c.JSON(http.StatusNotAcceptable, gin.H{"error": "Not Acceptable", "msg": err.Error()})
			return
		}
		submission.Code = string(codeBytes)

		res, err := k8s.K8sMgr.RunCodeTests(submission)

		if err != nil {
			res.Error = err.Error()
		}

		if res.Status != currentstatus.SUCCESS.ToString() {
			c.JSON(http.StatusNotAcceptable, res)
			return
		}
		if errors.Is(err, k8s.ErrTooManyRequests) {
			c.JSON(http.StatusTooManyRequests, res)
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

func PrivateTestSubmission() gin.HandlerFunc {
	return func(c *gin.Context) {
		var submission models.Submission
		if err := c.ShouldBind(&submission); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request", "msg": err.Error()})
			return
		}

		if len(submission.Tests) == 0 {
			c.JSON(http.StatusNoContent, gin.H{"error": "No tests provided"})
			return
		}
		codeBytes, err := base64.StdEncoding.DecodeString(submission.Code)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request", "msg": "The provided code is not properly base64 encoded."})
			return
		}
		if err := models.ValidateSubmission(&submission); err != nil {
			c.JSON(http.StatusNotAcceptable, gin.H{"error": "Not Acceptable", "msg": err.Error()})
			return
		}
		submission.Code = string(codeBytes)

		res, err := k8s.K8sMgr.RunCodeTests(submission)

		if err != nil {
			res.Error = err.Error()
		}
		if res.Status != currentstatus.SUCCESS.ToString() {
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

func RunCode(upgrader *websocket.Upgrader) gin.HandlerFunc {
	return func(c *gin.Context) {
		ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error", "msg": "Unable to upgrade connection to web socket"})
			return
		}
		defer ws.Close()
		var submission models.Submission
		if err := ws.ReadJSON(&submission); err != nil {
			ws.WriteMessage(websocket.TextMessage, []byte(err.Error()))
			return
		}

		if err := models.ValidateSubmission(&submission); err != nil {
			ws.WriteMessage(websocket.TextMessage, []byte(err.Error()))
			return
		}

		codeBytes, err := base64.StdEncoding.DecodeString(submission.Code)

		if err != nil {
			ws.WriteMessage(websocket.TextMessage, []byte("The code is not properly base64 encoded"))
			return
		}

		submission.Code = string(codeBytes)

		err = k8s.K8sMgr.RunCode(ws, submission)

		if err != nil {
			return
		}

	}
}
