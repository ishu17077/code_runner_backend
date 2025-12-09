package rust

import (
	"context"
	"fmt"
	"os/exec"
	"time"

	"github.com/ishu17077/code_runner_backend/worker-node/models"
	currentstatus "github.com/ishu17077/code_runner_backend/worker-node/models/enums/current_status"
	coderunners "github.com/ishu17077/code_runner_backend/worker-node/runner/helpers/code_runners"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func PreCompilationTask(submission models.Submission) (string, string, error) {
	newId := bson.NewObjectID().Hex()
	var dirPath = fmt.Sprintf("/temp/%s", newId)
	var filePath = fmt.Sprintf("%s/rustprog.rs", dirPath)
	var outputPath = fmt.Sprintf("%s/rustprog", dirPath)

	if err := coderunners.SaveFile(filePath, dirPath, submission.Code); err != nil {
		return "", dirPath, err
	}
	if err := compileCode(filePath, outputPath); err != nil {
		return filePath, dirPath, err
	}
	return outputPath, dirPath, nil
}

func CheckSubmission(submission models.Submission, test models.TestCase, binaryFilePath string) (currentstatus.CurrentStatus, error) {
	res, err := executeCode(binaryFilePath, test.Stdin)
	if err != nil {
		return currentstatus.FAILED, err
	}
	return coderunners.CheckOutput(res, test.ExpectedOutput)
}

func compileCode(filePath, outputPath string) error {
	var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "rustc", filePath, "-o", outputPath)

	coderunners.SetPermissions(cmd)
	res, err := cmd.CombinedOutput()
	coderunners.SetResourceLimits(cmd)

	if err != nil {
		return fmt.Errorf("Compilation Failed: %s %s", err.Error(), string(res))
	}
	return nil
}

func executeCode(binaryFilePath, stdin string) (string, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	runCmd := exec.CommandContext(ctx, binaryFilePath)
	return coderunners.RunCommandWithInput(runCmd, stdin)

}
