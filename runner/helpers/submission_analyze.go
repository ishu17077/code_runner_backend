package helpers

import (
	"fmt"
	"time"

	"github.com/ishu17077/code_runner_backend/models"
	currentstatus "github.com/ishu17077/code_runner_backend/models/enums/current_status"
	"github.com/ishu17077/code_runner_backend/models/enums/language"
	coderunners "github.com/ishu17077/code_runner_backend/runner/helpers/code_runners"
	"github.com/ishu17077/code_runner_backend/runner/helpers/code_runners/c"
	"github.com/ishu17077/code_runner_backend/runner/helpers/code_runners/c_sharp"
	"github.com/ishu17077/code_runner_backend/runner/helpers/code_runners/cpp"
	"github.com/ishu17077/code_runner_backend/runner/helpers/code_runners/golang"
	"github.com/ishu17077/code_runner_backend/runner/helpers/code_runners/java"
	"github.com/ishu17077/code_runner_backend/runner/helpers/code_runners/python"
	"github.com/ishu17077/code_runner_backend/runner/helpers/code_runners/rust"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func AnalyzeSubmission(submission models.Submission) (bool, []models.ExecResult, error) {
	lang := language.LanguageParser(submission.Language)
	var execResults []models.ExecResult
	switch lang {
	case language.C:
		res, err := testCCode(submission, &execResults)
		if err != nil || !res {
			return false, execResults, err
		}
		return true, execResults, err
	case language.Cpp:
		res, err := testCppCode(submission, &execResults)
		if err != nil || !res {
			return false, execResults, err
		}
		return true, execResults, err

	case language.Python:
		res, err := testPythonCode(submission, &execResults)
		if err != nil || !res {
			return false, execResults, err
		}
		return true, execResults, err

	case language.Java:
		res, err := testJavaCode(submission, &execResults)
		if err != nil || !res {
			return false, execResults, err
		}
		return true, execResults, err

	case language.Cs:
		res, err := testCSharpCode(submission, &execResults)
		if err != nil || !res {
			return false, execResults, err
		}
		return true, execResults, err
	case language.Rust:
		res, err := testRustCode(submission, &execResults)
		if err != nil || !res {
			return false, execResults, err
		}
		return true, execResults, err
	case language.Go:
		res, err := testGoCode(submission, &execResults)
		if err != nil || !res {
			return false, execResults, err
		}
		return true, execResults, err
	case language.Undefined:
		return false, []models.ExecResult{}, fmt.Errorf("Invalid Language Provided")
	}
	return false, execResults, nil
}

func testCCode(submission models.Submission, execResults *[]models.ExecResult) (bool, error) {
	var allPassed = true
	outputPath, dirPath, err := c.PreCompilationTask(submission)
	defer cleanUp(dirPath)
	if err != nil {
		return false, fmt.Errorf("Error compiling the file: %s", err.Error())
	}
	errors := 0
	for _, testCase := range submission.Tests {
		res, err := c.CheckSubmission(submission, testCase, outputPath)
		var execResult models.ExecResult
		execResult, passed := getExecResult(submission, testCase, res, err)
		if !passed {
			allPassed = false
		}
		*execResults = append(*execResults, execResult)
		if err != nil {
			errors++
			if errors > 2 {
				return false, nil
			}
		} else {
			errors = 0
		}

	}

	return allPassed, nil
}

func testCppCode(submission models.Submission, execResults *[]models.ExecResult) (bool, error) {
	var allPassed = true
	outputPath, dirPath, err := cpp.PreCompilationTask(submission)
	defer cleanUp(dirPath)
	if err != nil {
		return false, fmt.Errorf("Error compiling the file: %s", err.Error())
	}
	errors := 0
	for _, testCase := range submission.Tests {
		res, err := cpp.CheckSubmission(submission, testCase, outputPath)
		var execResult models.ExecResult
		execResult, passed := getExecResult(submission, testCase, res, err)
		if !passed {
			allPassed = false
		}
		*execResults = append(*execResults, execResult)

		if err != nil {
			errors++
			if errors > 2 {
				return false, nil
			}
		} else {
			errors = 0
		}
	}
	return allPassed, nil
}

func testPythonCode(submission models.Submission, execResults *[]models.ExecResult) (bool, error) {
	var allPassed = true
	filePath, dirPath, err := python.PreCompilationTask(submission)
	defer cleanUp(dirPath)
	if err != nil {
		return false, fmt.Errorf("Error compiling the file: %s", err.Error())
	}
	errors := 0
	for _, testCase := range submission.Tests {
		res, err := python.CheckSubmission(submission, testCase, filePath)
		var execResult models.ExecResult
		execResult, passed := getExecResult(submission, testCase, res, err)
		if !passed {
			allPassed = false
		}
		*execResults = append(*execResults, execResult)
		if err != nil {
			errors++
			if errors > 2 {
				return false, nil
			}
		} else {
			errors = 0
		}
	}
	return allPassed, nil
}

func testJavaCode(submission models.Submission, execResults *[]models.ExecResult) (bool, error) {
	var allPassed = false
	className, classPath, dirPath, err := java.PreCompilationTask(submission)
	defer cleanUp(dirPath)
	if err != nil {
		return false, fmt.Errorf("Error compiling the file: %s", err.Error())
	}
	var payload models.Payload = models.Payload{
		Class_name: "Solution",
		Exec_time:  2500,
		Tests:      submission.Tests,
	}
	allPassed, result, err := java.CheckSubmission(payload, className, classPath)
	if err != nil {
		allPassed = false

	}

	*execResults = result.Results

	return allPassed, err
}

func testCSharpCode(submission models.Submission, execResults *[]models.ExecResult) (bool, error) {
	var allPassed = true
	filePath, dirPath, err := c_sharp.PreCompilationTask(submission)
	defer cleanUp(dirPath)
	if err != nil {
		return false, fmt.Errorf("Error compiling the file: %s", err.Error())
	}
	errors := 0
	for _, testCase := range submission.Tests {
		res, err := c_sharp.CheckSubmission(submission, testCase, filePath)
		var execResult models.ExecResult
		execResult, passed := getExecResult(submission, testCase, res, err)
		if !passed {
			allPassed = false
		}
		*execResults = append(*execResults, execResult)
		if err != nil {
			errors++
			if errors > 2 {
				return false, nil
			}
		} else {
			errors = 0
		}
	}
	return allPassed, nil
}

func testRustCode(submission models.Submission, execResults *[]models.ExecResult) (bool, error) {
	var allPassed = true
	filePath, dirPath, err := rust.PreCompilationTask(submission)
	defer cleanUp(dirPath)
	if err != nil {
		return false, fmt.Errorf("Error compiling the file: %s", err.Error())
	}
	errors := 0
	for _, testCase := range submission.Tests {
		res, err := rust.CheckSubmission(submission, testCase, filePath)
		var execResult models.ExecResult
		execResult, passed := getExecResult(submission, testCase, res, err)
		if !passed {
			allPassed = false
		}
		*execResults = append(*execResults, execResult)

		if err != nil {
			errors++
			if errors > 2 {
				return false, nil
			}
		} else {
			errors = 0
		}
	}
	return allPassed, nil
}

func testGoCode(submission models.Submission, execResults *[]models.ExecResult) (bool, error) {
	var allPassed = true
	filepath, dirPath, err := golang.PreCompilationTask(submission)
	defer cleanUp(dirPath)
	if err != nil {
		return false, fmt.Errorf("Error compiling the file: %s", err.Error())
	}
	errors := 0
	for _, testCase := range submission.Tests {
		res, err := golang.CheckSubmission(submission, testCase, filepath)
		execResult, passed := getExecResult(submission, testCase, res, err)
		if !passed {
			allPassed = false
		}
		*execResults = append(*execResults, execResult)

		if err != nil {
			errors++
			if errors > 2 {
				return false, nil
			}
		} else {
			errors = 0
		}
	}
	return allPassed, nil
}

func getExecResult(submission models.Submission, testCase models.TestCase, res currentstatus.CurrentStatus, err error) (models.ExecResult, bool) {
	var execResult models.ExecResult = models.ExecResult{
		ExecResult_id: bson.NewObjectID().Hex(),
		Problem_id:    submission.ProblemID,
		Test_id:       testCase.Test_id,
	}
	if err != nil || res != currentstatus.SUCCESS {
		execResult.Status = &models.Status{
			Message:        err.Error(),
			Current_status: currentstatus.FAILED.ToString(),
			Stdout:         "",
			Stderr:         "",
			Completed_At:   time.Now(),
		}
		return execResult, false
	} else {
		//? SUCCESSful Execution :)
		execResult.Status = &models.Status{
			Message:        fmt.Sprintf("Test: #%s Passed", testCase.Test_id),
			Current_status: currentstatus.SUCCESS.ToString(),
			Stdout:         "",
			Stderr:         "",
			Completed_At:   time.Now(),
		}
		return execResult, true
	}

}

func cleanUp(path string) {
	go func(path string) {
		coderunners.CleanUp(path)
	}(path)
}
