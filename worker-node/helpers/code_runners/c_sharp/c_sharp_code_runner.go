package c_sharp

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

const filePath = "/temp/prog.cs"

func PreCompilationTask(submission models.Submission) error {
	if err := coderunners.SaveFile(filePath, submission.Code); err != nil {
		return err
	}
	return nil
}

func CheckSubmission(submission models.Submission, test models.TestCase) (currentstatus.CurrentStatus, error) {
	res, err := executeCode(filePath, test.Stdin)
	if err != nil {
		return currentstatus.FAILED, err
	}
	if strings.TrimSpace(res) == strings.TrimSpace(test.ExpectedOutput) {
		return currentstatus.SUCCESS, nil
	}
	return currentstatus.FAILED, fmt.Errorf("FAILED: Expected output: %s. Actual output: %s", test.ExpectedOutput, res)
}

func executeCode(filepath, stdin string) (string, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	runCmd := exec.CommandContext(ctx, "dotnet", "run", filepath)
	return coderunners.RunCommandWithInput(runCmd, stdin)
}
