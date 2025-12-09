package java

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

// const outputPath = "/temp/Solution.class"

func PreCompilationTask(submission models.Submission) (string, string, error) {
	newId := bson.NewObjectID().Hex()
	var dirPath = fmt.Sprintf("/temp/%s", newId)
	var filePath = fmt.Sprintf("%s/Solution.java", dirPath)
	const className = "Solution"

	if err := coderunners.SaveFile(filePath, dirPath, submission.Code); err != nil {
		return "", dirPath, err
	}
	if err := compileCode(filePath, dirPath); err != nil {
		return "", dirPath, err
	}
	return className, dirPath, nil
}

func CheckSubmission(submission models.Submission, test models.TestCase, className, dirPath string) (currentstatus.CurrentStatus, error) {
	res, err := executeCode(dirPath, className, test.Stdin)
	if err != nil {
		return currentstatus.FAILED, err
	}
	return coderunners.CheckOutput(res, test.ExpectedOutput)
}

func compileCode(filepath, outputDir string) error {
	var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "javac", "-d", outputDir, filepath)

	coderunners.SetPermissions(cmd)
	res, err := cmd.CombinedOutput()
	coderunners.SetResourceLimits(cmd)

	if err != nil {
		return fmt.Errorf("Compilation Failed: %s %s", err.Error(), string(res))
	}
	return nil
}

func executeCode(classPath, className, stdin string) (string, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()
	runCmd := exec.CommandContext(ctx, "java", "-cp", classPath, className)

	return coderunners.RunCommandWithInput(runCmd, stdin)
}
