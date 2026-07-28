package c

import (
	"context"
	"fmt"
	"os/exec"
	"time"

	"github.com/ishu17077/code_runner_backend/models"
	coderunners "github.com/ishu17077/code_runner_backend/runner/helpers/code_runners"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func PreCompilationTask(submission models.Submission) (string, string, error) {
	newId := bson.NewObjectID().Hex()
	var dirPath = fmt.Sprintf("/temp/%s", newId)
	var filePath = fmt.Sprintf("%s/main.c", dirPath)
	var outputPath = fmt.Sprintf("%s/main", dirPath)

	if err := coderunners.SaveFile(filePath, dirPath, submission.Code); err != nil {
		return "", dirPath, err
	}
	if err := compileCode(filePath, outputPath); err != nil {
		return filePath, dirPath, err
	}
	return outputPath, dirPath, nil
}

func ExecuteSubmission(test models.TestCase, binaryFile string) (string, error) {
	//TODO: Impl executeCcode test case
	return executeCode(binaryFile, test.Stdin)

}

func PipeSubmission(binaryFilePath string) error {
	var ctx, cancel = context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()
	runCmd := exec.CommandContext(ctx, binaryFilePath)
	return coderunners.PipeCommand(runCmd)
}

func compileCode(filePath string, outputPath string) error {
	var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "gcc", "-O2", "-include", "code_runner_backend_c_std.h", filePath, "-o", outputPath, "-lm")

	coderunners.SetPermissions(cmd)
	res, err := cmd.CombinedOutput()

	if err != nil {
		return fmt.Errorf("Compilation Failed: %s %s", err.Error(), string(res))
	}
	return nil
}

func executeCode(binaryFilePath string, stdin string) (string, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	runCmd := exec.CommandContext(ctx, binaryFilePath)
	return coderunners.RunCommandWithInput(runCmd, ctx, stdin)
}
