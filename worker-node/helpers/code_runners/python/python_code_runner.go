package python

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

const filePath = "/temp/main.py"

func PreCompilationTask(submission models.Submission) error {
	if err := coderunners.SaveFile(filePath, submission.Code); err != nil {
		return fmt.Errorf("Error saving the file: %s", err.Error())
	}
	return nil
}

func CheckSubmission(submission models.Submission, test models.TestCase) (currentstatus.CurrentStatus, error) {
	res, err := executeCode(filePath, test.Stdin)
	if err != nil {
		return currentstatus.FAILED, fmt.Errorf("The test was unsuccessful: %s", err.Error())
	}

	if strings.TrimSpace(res) == strings.TrimSpace(test.ExpectedOutput) {
		return currentstatus.SUCCESS, nil
	}
	return currentstatus.SUCCESS, fmt.Errorf("Test: #%s Failed", test.Test_id)
}

func executeCode(filePath string, stdin string) (string, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	runCmd := exec.CommandContext(ctx, "python", filePath)
	return coderunners.RunCommandWithInput(runCmd, stdin)
}
