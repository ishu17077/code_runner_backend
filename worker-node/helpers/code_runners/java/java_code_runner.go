package java

import (
	"context"
	"fmt"
	"os/exec"
	"time"

	coderunners "github.com/ishu17077/code_runner_backend/worker-node/helpers/code_runners"
	"github.com/ishu17077/code_runner_backend/worker-node/models"
	currentstatus "github.com/ishu17077/code_runner_backend/worker-node/models/enums/current_status"
)

const filePath = "/temp/Solution.java"

// const outputPath = "/temp/Solution.class"
const className = "Solution"
const dir = "/temp/"

func PreCompilationTask(submission models.Submission) error {
	if err := coderunners.SaveFile(filePath, submission.Code); err != nil {
		return err
	}

	if err := compileCode(filePath, dir); err != nil {
		return err
	}
	return nil
}

func CheckSubmission(submission models.Submission, test models.TestCase) (currentstatus.CurrentStatus, error) {
	res, err := executeCode(dir, className, test.Stdin)
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
