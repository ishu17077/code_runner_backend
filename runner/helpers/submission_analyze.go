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
	"github.com/ishu17077/code_runner_backend/runner/helpers/code_runners/javascript"
	"github.com/ishu17077/code_runner_backend/runner/helpers/code_runners/perl"
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
	case language.Perl:
		res, err := testPerlCode(submission, &execResults)
		if err != nil || !res {
			return false, execResults, err
		}
		return true, execResults, err
	case language.Javascript:
		res, err := testJavascriptCode(submission, &execResults)
		if err != nil || !res {
			return false, execResults, err
		}
		return true, execResults, err
	case language.Undefined:
		return false, []models.ExecResult{}, fmt.Errorf("INVALID_LANGUAGE_PROVIDED")
	}
	return false, execResults, nil
}

func testCCode(submission models.Submission, execResults *[]models.ExecResult) (bool, error) {
	var allPassed = true
	var gErr error = nil
	outputPath, dirPath, err := c.PreCompilationTask(submission)
	defer cleanUp(dirPath)
	if err != nil {
		return false, fmt.Errorf("Error compiling the file: %s", err.Error())
	}
	errors := 0
	for _, testCase := range submission.Tests {
		startTime := time.Now()
		output, err := c.ExecuteSubmission(testCase, outputPath)
		execTime := time.Since(startTime).Milliseconds()
		execResult, passed := getExecResult(testCase, output, err)
		execResult.Status.Exec_time_ms = uint32(execTime)
		if !passed {
			allPassed = false
			gErr = fmt.Errorf("%s", currentstatus.FAILED.ToString())
		}
		*execResults = append(*execResults, execResult)
		if err != nil {
			errors++
			if errors > 2 {
				return false, err
			}
		} else {
			errors = 0
		}

	}

	return allPassed, gErr
}

func testCppCode(submission models.Submission, execResults *[]models.ExecResult) (bool, error) {
	var allPassed = true
	var gErr error = nil
	outputPath, dirPath, err := cpp.PreCompilationTask(submission)
	defer cleanUp(dirPath)
	if err != nil {
		return false, fmt.Errorf("Error compiling the file: %s", err.Error())
	}
	errors := 0
	for _, testCase := range submission.Tests {
		startTime := time.Now()
		output, err := cpp.ExecuteSubmission(testCase, outputPath)
		execTime := time.Since(startTime).Milliseconds()
		execResult, passed := getExecResult(testCase, output, err)
		execResult.Status.Exec_time_ms = uint32(execTime)
		if !passed {
			allPassed = false
			gErr = fmt.Errorf("%s", currentstatus.FAILED.ToString())
		}
		*execResults = append(*execResults, execResult)

		if err != nil {
			errors++
			if errors > 2 {
				return false, err
			}
		} else {
			errors = 0
		}
	}
	return allPassed, gErr
}

func testPythonCode(submission models.Submission, execResults *[]models.ExecResult) (bool, error) {
	var allPassed = true
	var gErr error = nil
	filePath, dirPath, err := python.PreCompilationTask(submission)
	defer cleanUp(dirPath)
	if err != nil {
		return false, fmt.Errorf("Error processing the file: %s", err.Error())
	}
	errors := 0
	for _, testCase := range submission.Tests {
		startTime := time.Now()
		output, err := python.ExecuteSubmission(testCase, filePath)
		execTime := time.Since(startTime).Milliseconds()
		execResult, passed := getExecResult(testCase, output, err)
		execResult.Status.Exec_time_ms = uint32(execTime)
		if !passed {
			allPassed = false
			gErr = fmt.Errorf("%s", currentstatus.FAILED.ToString())
		}
		*execResults = append(*execResults, execResult)
		if err != nil {
			errors++
			if errors > 2 {
				return false, err
			}
		} else {
			errors = 0
		}
	}
	return allPassed, gErr
}

func testJavaCode(submission models.Submission, execResults *[]models.ExecResult) (bool, error) {
	var allPassed = false
	className, classPath, dirPath, err := java.PreCompilationTask(submission)
	defer cleanUp(dirPath)
	if err != nil {
		return false, fmt.Errorf("Error compiling the file: %s", err.Error())
	}
	var payload models.JavaDriverPayload = models.JavaDriverPayload{
		Class_name: "Solution",
		Exec_time:  2500,
		Tests:      submission.Tests,
	}
	allPassed, result, err := java.ExecuteSubmission(payload, className, classPath)
	if err != nil {
		allPassed = false
	}

	*execResults = result.Results

	return allPassed, err
}

func testCSharpCode(submission models.Submission, execResults *[]models.ExecResult) (bool, error) {
	var allPassed = true
	var gErr error
	filePath, dirPath, err := c_sharp.PreCompilationTask(submission)
	defer cleanUp(dirPath)
	if err != nil {
		return false, fmt.Errorf("Error compiling the file: %s", err.Error())
	}
	errors := 0
	for _, testCase := range submission.Tests {
		startTime := time.Now()
		output, err := c_sharp.ExecuteSubmission(testCase, filePath)
		execTime := time.Since(startTime).Milliseconds()
		execResult, passed := getExecResult(testCase, output, err)
		execResult.Status.Exec_time_ms = uint32(execTime)
		if !passed {
			allPassed = false
			gErr = fmt.Errorf("%s", currentstatus.FAILED.ToString())
		}
		*execResults = append(*execResults, execResult)
		if err != nil {
			errors++
			if errors > 2 {
				return false, err
			}
		} else {
			errors = 0
		}
	}
	return allPassed, gErr
}

func testRustCode(submission models.Submission, execResults *[]models.ExecResult) (bool, error) {
	var allPassed = true
	var gErr error = nil
	outputPath, dirPath, err := rust.PreCompilationTask(submission)
	defer cleanUp(dirPath)
	if err != nil {
		return false, fmt.Errorf("Error compiling the file: %s", err.Error())
	}
	errors := 0
	for _, testCase := range submission.Tests {
		startTime := time.Now()
		output, err := rust.ExecuteSubmission(testCase, outputPath)
		execTime := time.Since(startTime).Milliseconds()
		execResult, passed := getExecResult(testCase, output, err)
		execResult.Status.Exec_time_ms = uint32(execTime)
		if !passed {
			allPassed = false
			gErr = fmt.Errorf("%s", currentstatus.FAILED.ToString())
		}
		*execResults = append(*execResults, execResult)

		if err != nil {
			errors++
			if errors > 2 {
				return false, err
			}
		} else {
			errors = 0
		}
	}
	return allPassed, gErr
}

func testGoCode(submission models.Submission, execResults *[]models.ExecResult) (bool, error) {
	var allPassed = true
	var gErr error = nil
	outputPath, dirPath, err := golang.PreCompilationTask(submission)
	defer cleanUp(dirPath)
	if err != nil {
		return false, fmt.Errorf("Error compiling the file: %s", err.Error())
	}
	errors := 0
	for _, testCase := range submission.Tests {
		startTime := time.Now()
		output, err := golang.ExecuteSubmission(testCase, outputPath)
		execTime := time.Since(startTime).Milliseconds()
		execResult, passed := getExecResult(testCase, output, err)
		execResult.Status.Exec_time_ms = uint32(execTime)
		if !passed {
			allPassed = false
			gErr = fmt.Errorf("%s", currentstatus.FAILED.ToString())
		}
		*execResults = append(*execResults, execResult)

		if err != nil {
			errors++
			if errors > 2 {
				return false, err
			}
		} else {
			errors = 0
		}
	}
	return allPassed, gErr
}

func testPerlCode(submission models.Submission, execResults *[]models.ExecResult) (bool, error) {
	var allPassed = true
	var gErr error = nil
	filePath, dirPath, err := perl.PreCompilationTask(submission)
	defer cleanUp(dirPath)

	if err != nil {
		return false, fmt.Errorf("Error processing the file: %s", err.Error())
	}

	errors := 0
	for _, testCase := range submission.Tests {
		startTime := time.Now()
		output, err := perl.ExecuteSubmission(testCase, filePath)
		execTime := time.Since(startTime).Milliseconds()
		execResult, passed := getExecResult(testCase, output, err)
		execResult.Status.Exec_time_ms = uint32(execTime)
		if !passed {
			allPassed = false
			gErr = fmt.Errorf("%s", currentstatus.FAILED.ToString())
		}
		*execResults = append(*execResults, execResult)
		if err != nil {
			errors++
			if errors > 2 {
				return false, err
			}
		} else {
			errors = 0
		}
	}
	return allPassed, gErr
}

func testJavascriptCode(submission models.Submission, execResults *[]models.ExecResult) (bool, error) {
	var allPassed = true
	var gErr error = nil
	filePath, dirPath, err := javascript.PreCompilationTask(submission)
	defer cleanUp(dirPath)

	if err != nil {
		return false, fmt.Errorf("Error processing the file: %s", err.Error())
	}
	errors := 0
	for _, testCase := range submission.Tests {
		startTime := time.Now()
		output, err := javascript.ExecuteSubmission(testCase, filePath)
		execTime := time.Since(startTime).Milliseconds()
		execResult, passed := getExecResult(testCase, output, err)
		execResult.Status.Exec_time_ms = uint32(execTime)
		if !passed {
			allPassed = false
			gErr = fmt.Errorf("%s", currentstatus.FAILED.ToString())
		}
		*execResults = append(*execResults, execResult)
		if err != nil {
			errors++
			if errors > 2 {
				return false, err
			}
		} else {
			errors = 0
		}
	}
	return allPassed, gErr
}

func getExecResult(testCase models.TestCase, output string, err error) (models.ExecResult, bool) {
	var execResult models.ExecResult = models.ExecResult{
		ExecResult_id: bson.NewObjectID().Hex(),
		Test_id:       testCase.Test_id,
	}

	if err != nil {
		if err != coderunners.TleError {
			execResult.Status = &models.Status{
				Message:         err.Error(),
				Current_status:  currentstatus.FAILED.ToString(),
				Stdout:          output,
				Stdin:           testCase.Stdin,
				Expected_output: testCase.ExpectedOutput,
				Stderr:          err.Error(),
				Completed_At:    time.Now(),
			}
			return execResult, false
		} else {
			execResult.Status = &models.Status{
				Message:         err.Error(),
				Current_status:  currentstatus.TIME_LIMIT_EXCEEDED.ToString(),
				Stdout:          output,
				Stdin:           testCase.Stdin,
				Expected_output: testCase.ExpectedOutput,
				Stderr:          "",
				Completed_At:    time.Now(),
			}
		}
		return execResult, false
	}
	var status, testResErr = coderunners.CheckOutput(output, testCase.ExpectedOutput)
	if testResErr != nil || status != currentstatus.SUCCESS {
		execResult.Status = &models.Status{
			Message:         testResErr.Error(),
			Current_status:  currentstatus.FAILED.ToString(),
			Stdout:          output,
			Stdin:           testCase.Stdin,
			Expected_output: testCase.ExpectedOutput,
			Stderr:          "",
			Completed_At:    time.Now(),
		}
		return execResult, false
	} else {
		//? SUCCESSful Execution :)
		execResult.Status = &models.Status{
			Message:         fmt.Sprintf("Test: #%s Passed", testCase.Test_id),
			Current_status:  currentstatus.SUCCESS.ToString(),
			Stdout:          output,
			Stderr:          "",
			Stdin:           testCase.Stdin,
			Expected_output: testCase.ExpectedOutput,
			Completed_At:    time.Now(),
		}
		return execResult, true
	}

}

func cleanUp(dirPath string) {
	go func(path string) {
		coderunners.CleanUp(path)
	}(dirPath)
}
