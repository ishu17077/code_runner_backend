package c

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"

	coderunners "github.com/ishu17077/code_runner_backend/worker-node/helpers/code_runners"
	"github.com/ishu17077/code_runner_backend/worker-node/models"
	currentstatus "github.com/ishu17077/code_runner_backend/worker-node/models/enums/current_status"
)

const filePath = "/temp/main.c"
const outputPath = "/temp/main"

func PreCompilationTask(submission models.Submission) error {
	if err := coderunners.SaveFile(filePath, submission.Code); err != nil {
		return err
	}

	if err := compileCode(filePath, outputPath); err != nil {
		return fmt.Errorf("Compilation Failed: %s", err.Error())
	}
	return nil
}

func CheckSubmission(submission models.Submission, test models.TestCase) (currentstatus.CurrentStatus, error) {

	//TODO: Impl executeCcode test case
	res, err := executeCode(outputPath, test.Stdin)
	if err != nil {
		return currentstatus.FAILED, err
	}

	if strings.TrimSpace(res) == strings.TrimSpace(test.ExpectedOutput) {
		return currentstatus.SUCCESS, nil
	}
	return currentstatus.FAILED, fmt.Errorf("FAILED: Expected output: %s. Actual output: %s", test.ExpectedOutput, res)
}

func compileCode(filePath string, outputPath string) error {
	cmd := exec.Command("gcc", filePath, "-o", outputPath, "-lm")
	_, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("Compilation Failed: %s", err.Error())
	}
	// fileMode := os.FileMode(0755)
	// if chmodErr := os.Chmod("/temp/main", fileMode); chmodErr != nil {
	// 	return fmt.Errorf("Failed to set execute permissions to file")
	// }
	return nil
}

func executeCode(binaryFilePath string, stdin string) (string, error) {

	var ctx, cancel = context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	runCmd := exec.CommandContext(ctx, binaryFilePath)
	return coderunners.RunCommandWithInput(runCmd, stdin)
}
