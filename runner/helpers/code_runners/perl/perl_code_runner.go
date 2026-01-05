package perl

import (
	"fmt"
	"os/exec"

	"github.com/ishu17077/code_runner_backend/models"
	currentstatus "github.com/ishu17077/code_runner_backend/models/enums/current_status"
	coderunners "github.com/ishu17077/code_runner_backend/runner/helpers/code_runners"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func PreCompilationTask(submission models.Submission) (string, string, error) {
	newId := bson.NewObjectID().Hex()
	var dirPath = fmt.Sprintf("/temp/%s", newId)
	var filePath = fmt.Sprintf("%s/main.pl", dirPath)

	if err := coderunners.SaveFile(filePath, dirPath, submission.Code); err != nil {
		return "", dirPath, err
	}
	return filePath, dirPath, nil
}

func CheckSubmission(testcase models.TestCase, filePath string) (currentstatus.CurrentStatus, error) {
	res, err := executeCode(filePath, testcase.Stdin)
	if err != nil {
		return currentstatus.FAILED, err
	}
	return coderunners.CheckOutput(res, testcase.ExpectedOutput)
}

func executeCode(filePath string, stdin string) (string, error) {
	runCmd := exec.Command("perl", filePath)
	return coderunners.RunCommandWithInput(runCmd, stdin)
}
