package helpers

import (
	"fmt"
	"time"

	coderunners "github.com/ishu17077/code_runner_backend/worker-node/helpers/code_runners"
	"github.com/ishu17077/code_runner_backend/worker-node/helpers/code_runners/c"
	cs "github.com/ishu17077/code_runner_backend/worker-node/helpers/code_runners/c_sharp"
	"github.com/ishu17077/code_runner_backend/worker-node/helpers/code_runners/cpp"
	"github.com/ishu17077/code_runner_backend/worker-node/helpers/code_runners/java"
	"github.com/ishu17077/code_runner_backend/worker-node/helpers/code_runners/python"
	"github.com/ishu17077/code_runner_backend/worker-node/models"
	currentstatus "github.com/ishu17077/code_runner_backend/worker-node/models/enums/current_status"
	"github.com/ishu17077/code_runner_backend/worker-node/models/enums/language"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func AnalyzeSubmission(submission models.Submission, testCases []models.TestCase) (bool, []models.ExecResult, error) {
	lang := language.LanguageParser(submission.Language)
	var execResults []models.ExecResult
	switch lang {
	case language.C:
		res, err := testCCode(submission, testCases, &execResults)
		if err != nil || !res {
			return false, execResults, err
		}
		return true, execResults, err
	case language.Cpp:
		res, err := testCppCode(submission, testCases, &execResults)
		if err != nil || !res {
			return false, execResults, err
		}
		return true, execResults, err

	case language.Python:
		res, err := testPythonCode(submission, testCases, &execResults)
		if err != nil || !res {
			return false, execResults, err
		}
		return true, execResults, err

	case language.Java:
		res, err := testJavaCode(submission, testCases, &execResults)
		if err != nil || !res {
			return false, execResults, err
		}
		return true, execResults, err

	case language.Cs:
		res, err := testCSharpCode(submission, testCases, &execResults)
		if err != nil || !res {
			return false, execResults, err
		}
		return true, execResults, err
	}
	return false, execResults, nil
}

func testCCode(submission models.Submission, testCases []models.TestCase, execResults *[]models.ExecResult) (bool, error) {
	var allPassed = true
	defer cleanUp()
	if err := c.PreCompilationTask(submission); err != nil {
		return false, fmt.Errorf("Error compiling the file: %s", err.Error())
	}
	for _, testCase := range testCases {
		res, err := c.CheckSubmission(submission, testCase)
		var execResult models.ExecResult
		execResult, passed := getExecResults(submission, testCase, res, err)
		if !passed {
			allPassed = false
		}
		*execResults = append(*execResults, execResult)
	}

	return allPassed, nil
}

func testCppCode(submission models.Submission, testCases []models.TestCase, execResults *[]models.ExecResult) (bool, error) {
	var allPassed = true
	defer cleanUp()
	if err := cpp.PreCompilationTask(submission); err != nil {
		return false, fmt.Errorf("Error compiling the file: %s", err.Error())
	}
	for _, testCase := range testCases {
		res, err := cpp.CheckSubmission(submission, testCase)
		var execResult models.ExecResult
		execResult, passed := getExecResults(submission, testCase, res, err)
		if !passed {
			allPassed = false
		}
		*execResults = append(*execResults, execResult)
	}
	return allPassed, nil
}

func testPythonCode(submission models.Submission, testCases []models.TestCase, execResults *[]models.ExecResult) (bool, error) {
	var allPassed = true
	defer cleanUp()
	if err := python.PreCompilationTask(submission); err != nil {
		return false, fmt.Errorf("Error compiling the file: %s", err.Error())
	}
	for _, testCase := range testCases {
		res, err := python.CheckSubmission(submission, testCase)
		var execResult models.ExecResult
		execResult, passed := getExecResults(submission, testCase, res, err)
		if !passed {
			allPassed = false
		}
		*execResults = append(*execResults, execResult)
	}
	return allPassed, nil
}

func testJavaCode(submission models.Submission, testCases []models.TestCase, execResults *[]models.ExecResult) (bool, error) {
	var allPassed = true
	defer cleanUp()
	if err := java.PreCompilationTask(submission); err != nil {
		return false, fmt.Errorf("Error compiling the file: %s", err.Error())
	}
	for _, testCase := range testCases {
		res, err := java.CheckSubmission(submission, testCase)
		var execResult models.ExecResult
		execResult, passed := getExecResults(submission, testCase, res, err)
		if !passed {
			allPassed = false
		}
		*execResults = append(*execResults, execResult)
	}
	return allPassed, nil
}

func testCSharpCode(submission models.Submission, testCases []models.TestCase, execResults *[]models.ExecResult) (bool, error) {
	var allPassed = true
	defer cleanUp()
	if err := cs.PreCompilationTask(submission); err != nil {
		return false, fmt.Errorf("Error compiling the file: %s", err.Error())
	}
	for _, testCase := range testCases {
		res, err := cs.CheckSubmission(submission, testCase)
		var execResult models.ExecResult
		execResult, passed := getExecResults(submission, testCase, res, err)
		if !passed {
			allPassed = false
		}
		*execResults = append(*execResults, execResult)
	}
	return allPassed, nil
}

func getExecResults(submission models.Submission, testCase models.TestCase, res currentstatus.CurrentStatus, err error) (models.ExecResult, bool) {
	var execResult models.ExecResult = models.ExecResult{
		ID:         bson.NewObjectID(),
		Problem_id: submission.ProblemID,
		Team_id:    submission.Team_id,
		Test_id:    testCase.Test_id,
	}
	execResult.ExecResult_id = execResult.ID.Hex()
	if err != nil || res != currentstatus.SUCCESS {
		execResult.Status = &models.Status{
			Message:        err.Error(),
			Current_status: "FAILED",
			Stdout:         "",
			Stderr:         "",
			Completed_At:   time.Now(),
		}
		return execResult, false
	} else {
		//? SUCCESSful Execution :)
		execResult.Status = &models.Status{
			Message:        fmt.Sprintf("Test: #%s Passed", testCase.Test_id),
			Current_status: "SUCCESS",
			Stdout:         "",
			Stderr:         "",
			Completed_At:   time.Now(),
		}
		return execResult, true
	}

}

func cleanUp() {
	go func() {
		coderunners.CleanUp()
	}()
}
