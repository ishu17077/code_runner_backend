package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/ishu17077/code_runner_backend/worker-node/models"
	currentstatus "github.com/ishu17077/code_runner_backend/worker-node/models/enums/current_status"
	"github.com/ishu17077/code_runner_backend/worker-node/models/enums/language"
	"github.com/ishu17077/code_runner_backend/worker-node/runner/helpers"
)

var validate = validator.New()

func main() {
	inputBytes, err := io.ReadAll(os.Stdin)

	if err != nil {
		printInternalError("Failed to read stdin", err)
		return
	}
	rawInput := strings.TrimSpace(string(inputBytes))

	if len(rawInput) == 0 {
		printInternalError("Empty Input Received", err)
		return
	}

	inputBytes, err = base64.StdEncoding.DecodeString(rawInput)

	if err != nil {
		printInternalError("Invalid Base64 encoding provided", err)
		return
	}

	var submission models.Submission
	if err := json.Unmarshal(inputBytes, &submission); err != nil {
		printInternalError("Invalid JSON Payload", err)
		return
	}
	if validationErr := validate.Struct(submission); validationErr != nil {
		printInternalError("Important values absent from payload", validationErr)
		return
	}

	submission.Language = language.LanguageParser(submission.Language).ToString()

	allPassed, execResults, err := helpers.AnalyzeSubmission(submission, submission.Tests)
	var result models.Result = models.Result{
		Results: execResults,
	}
	if allPassed {
		result.Status = "SUCCESS"
	} else {
		result.Status = "FAILED"
	}

	if err != nil {
		result.Error = err.Error()
	}

	printFinalResult(result)
	// var result models.Result
	// result.Results = execResults
	// if allPassed {
	// 	result.Status = currentstatus.SUCCESS.ToString()
	// } else {
	// 	result.Status = currentstatus.FAILED.ToString()
	// }
	// printFinalResult(result)
}

func printFinalResult(res models.Result) {
	fmt.Println("---JSON_START---")
	jsonData, err := json.Marshal(res)
	if err != nil {
		fmt.Printf(err.Error())
		fmt.Println("---JSON_END---")
		return
	}
	base64enc := base64.StdEncoding.EncodeToString(jsonData)
	fmt.Println(base64enc)
	fmt.Print("---JSON_END---")
}

func printInternalError(msg string, err error) {
	fullMsg := msg
	if err != nil {
		fullMsg = fmt.Sprintf("%s: %v", msg, err)
	}

	res := models.Result{
		Status: "INTERNAL_ERROR",
		Results: []models.ExecResult{
			{
				Status: &models.Status{
					Message:        fullMsg,
					Current_status: currentstatus.FAILED.ToString(),
					Stdout:         "",
					Stderr:         "",
					Completed_At:   time.Now(),
				},
			}},
	}
	printFinalResult(res)
}
