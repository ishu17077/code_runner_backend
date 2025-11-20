package helpers

import (
	"fmt"
	"time"

	"github.com/ishu17077/code_runner_backend/worker-node/helpers/code_runners/c"
	"github.com/ishu17077/code_runner_backend/worker-node/helpers/code_runners/cpp"
	"github.com/ishu17077/code_runner_backend/worker-node/helpers/code_runners/python"
	"github.com/ishu17077/code_runner_backend/worker-node/models"
	"github.com/ishu17077/code_runner_backend/worker-node/models/enums"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func AnalyzeSubmission(submission models.Submission, testCases []models.TestCase) (bool, []models.ExecResult) {
	language := enums.LanguageParser(submission.Language)
	var execResults []models.ExecResult
	switch language {
	case enums.C:
		res, err := testCCode(submission, testCases, &execResults)
		if err != nil || !res {
			return false, execResults
		}
		return true, execResults
	case enums.Cpp:
		res, err := testCppCode(submission, testCases, &execResults)
		if err != nil || !res {
			return false, execResults
		}
		return true, execResults

	case enums.Python:
		res, err := testPythonCode(submission, testCases, &execResults)
		if err != nil || !res {
			return false, execResults
		}
		return true, execResults
	}
	return false, execResults
}

func testCCode(submission models.Submission, testCases []models.TestCase, execResults *[]models.ExecResult) (bool, error) {
	var allPassed = true
	if err := c.PreCompilationTask(submission); err != nil {
		return false, fmt.Errorf("Error Compiling the file the file")
	}
	for _, testCase := range testCases {
		res, err := c.CheckSubmission(submission, testCase)
		var execResult models.ExecResult
		execResult, passed := getExecResults(submission, testCase, err, res.ToString())
		if !passed {
			allPassed = false
		}
		*execResults = append(*execResults, execResult)
	}
	return allPassed, nil
}

func testCppCode(submission models.Submission, testCases []models.TestCase, execResults *[]models.ExecResult) (bool, error) {
	var allPassed = true
	if err := cpp.PreCompilationTask(submission); err != nil {
		return false, fmt.Errorf("Error Compiling the file the file")
	}
	for _, testCase := range testCases {
		res, err := cpp.CheckSubmission(submission, testCase)
		var execResult models.ExecResult
		execResult, passed := getExecResults(submission, testCase, err, res.ToString())
		if !passed {
			allPassed = false
		}
		*execResults = append(*execResults, execResult)
	}
	return allPassed, nil
}

func testPythonCode(submission models.Submission, testCases []models.TestCase, execResults *[]models.ExecResult) (bool, error) {
	var allPassed = true
	if err := python.PreCompilationTask(submission); err != nil {
		return false, fmt.Errorf("Error Compiling the file the file")
	}
	for _, testCase := range testCases {
		res, err := python.CheckSubmission(submission, testCase)
		var execResult models.ExecResult
		execResult, passed := getExecResults(submission, testCase, err, res.ToString())
		if !passed {
			allPassed = false
		}
		*execResults = append(*execResults, execResult)
	}
	return allPassed, nil
}

func getExecResults(submission models.Submission, testCase models.TestCase, err error, res string) (models.ExecResult, bool) {
	var execResult models.ExecResult = models.ExecResult{
		ID:         bson.NewObjectID(),
		Problem_id: submission.ProblemID,
		Team_id:    submission.Team_id,
		Test_id:    testCase.Test_id,
	}
	execResult.ExecResult_id = execResult.ID.Hex()
	if err != nil || res != "SUCCESS" {
		execResult.Status = &models.Status{
			Message:        err.Error(),
			Current_status: "FAILED",
			Stdout:         "",
			Stderr:         "",
			Completed_At:   time.Now(),
		}
		return execResult, false
	} else {
		//! SUCCESSful Execution
		execResult.Status = &models.Status{
			Message:        fmt.Sprintf("Test Case %s Passed", testCase.Test_id),
			Current_status: "SUCCESS",
			Stdout:         "",
			Stderr:         "",
			Completed_At:   time.Now(),
		}
		return execResult, true
	}

}
