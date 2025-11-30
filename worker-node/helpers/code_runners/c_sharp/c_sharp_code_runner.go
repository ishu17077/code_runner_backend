package c_sharp

import (
	"context"
	"os/exec"
	"time"

	coderunners "github.com/ishu17077/code_runner_backend/worker-node/helpers/code_runners"
	"github.com/ishu17077/code_runner_backend/worker-node/models"
	currentstatus "github.com/ishu17077/code_runner_backend/worker-node/models/enums/current_status"
)

//TODO: Compile with dotnet 10 sdk the 9 doens't really work

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
	return coderunners.CheckOutput(res, test.ExpectedOutput)
}

func executeCode(filepath, stdin string) (string, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	runCmd := exec.CommandContext(ctx, "dotnet-script", filepath)
	return coderunners.RunCommandWithInput(runCmd, stdin)
}
