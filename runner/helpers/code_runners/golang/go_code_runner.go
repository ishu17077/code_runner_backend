package golang

import (
	"context"
	"fmt"
	"os/exec"
	"time"

	"github.com/ishu17077/code_runner_backend/models"
	currentstatus "github.com/ishu17077/code_runner_backend/models/enums/current_status"
	coderunners "github.com/ishu17077/code_runner_backend/runner/helpers/code_runners"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func PreCompilationTask(submission models.Submission) (string, string, error) {
	newId := bson.NewObjectID().Hex()
	var dirPath = fmt.Sprintf("/temp/%s", newId)
	var filePath = fmt.Sprintf("%s/main.go", dirPath)
	var outputPath = fmt.Sprintf("%s/main", dirPath)

	if err := coderunners.SaveFile(filePath, dirPath, submission.Code); err != nil {
		return filePath, dirPath, err
	}

	if err := compileCode(filePath, outputPath); err != nil {
		return "", dirPath, err
	}

	return outputPath, dirPath, nil

}

func CheckSubmission(submission models.Submission, testCase models.TestCase, binaryFilePath string) (currentstatus.CurrentStatus, error) {
	res, err := executeCode(binaryFilePath, testCase.Stdin)
	if err != nil {
		return currentstatus.FAILED, err
	}
	return coderunners.CheckOutput(res, testCase.ExpectedOutput)
}

func compileCode(filePath, outputPath string) error {
	var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "go", "build", "-o", outputPath)
	coderunners.SetPermissions(cmd)

	_, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("Compilation Failed: %s", err.Error())
	}

	return nil
}

func executeCode(binaryFilePath, stdin string) (string, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	runCmd := exec.CommandContext(ctx, binaryFilePath)
	return coderunners.RunCommandWithInput(runCmd, stdin)
}
