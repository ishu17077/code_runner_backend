package python

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"time"

	coderunners "github.com/ishu17077/code_runner_backend/worker-node/helpers/code_runners"
	"github.com/ishu17077/code_runner_backend/worker-node/models"
	"github.com/ishu17077/code_runner_backend/worker-node/models/enums"
)

const filePath = "/temp/main.py"

func PreCompilationTask(submission models.Submission) error {
	if err := coderunners.SaveFile(filePath, submission.Code); err != nil {
		return fmt.Errorf("Error saving the file: %s", err.Error())
	}
	return nil
}

func CheckSubmission(submission models.Submission, test models.TestCase) (enums.CurrentStatus, error) {
	res, err := executePythonCode(filePath, test.Stdin)
	if err != nil {
		return enums.FAILED, fmt.Errorf("The test was unsuccessful: %s", err.Error())
	}

	if strings.TrimSpace(res) == strings.TrimSpace(test.ExpectedOutput) {
		return enums.SUCCESS, nil
	}
	return enums.SUCCESS, fmt.Errorf("Test: #%s Failed", test.Test_id)
}

func executePythonCode(filePath string, stdin string) (string, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	runCmd := exec.CommandContext(ctx, "python", filePath)
	// coderunners.SetLimitsAndPermissions(runCmd)
	stdinPipe, pipeErr := runCmd.StdinPipe()

	var outputBuffer bytes.Buffer

	runCmd.Stdout = &outputBuffer
	if pipeErr != nil {
		return "", fmt.Errorf("Error connecting pipe input")
	}
	if startErr := runCmd.Start(); startErr != nil {
		return "", fmt.Errorf("Error starting the program")
	}

	_, writeErr := io.WriteString(stdinPipe, stdin)
	if writeErr != nil {
		return "", fmt.Errorf("Error writing input to file")
	}

	stdinPipe.Close()

	if waitErr := runCmd.Wait(); waitErr != nil {
		return "", fmt.Errorf("Error executing file")
	}

	var finalOutput = outputBuffer.String()

	return finalOutput, nil
}
