package cpp

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
	currentstatus "github.com/ishu17077/code_runner_backend/worker-node/models/enums/current_status"
)

const filePath = "/temp/main.cpp"
const outputPath = "/temp/main"

func PreCompilationTask(submission models.Submission) error {

	if err := coderunners.SaveFile(filePath, submission.Code); err != nil {
		return err
	}

	if err := compileCode(filePath, outputPath); err != nil {
		return err
	}

	return nil
}

func CheckSubmission(submission models.Submission, test models.TestCase) (currentstatus.CurrentStatus, error) {
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
	cmd := exec.Command("g++", filePath, "-o", outputPath)

	_, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("Compilation Failed: %s", err.Error())
	}

	return nil
}

func executeCode(binaryFilePath string, stdin string) (string, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	runCmd := exec.CommandContext(ctx, binaryFilePath)
	coderunners.SetPermissions(runCmd)
	stdinPipe, pipErr := runCmd.StdinPipe()

	if pipErr != nil {
		return "", fmt.Errorf("Error connecting pipe input")
	}

	var outputBuffer bytes.Buffer

	runCmd.Stdout = &outputBuffer

	if startErr := runCmd.Start(); startErr != nil {
		return "", fmt.Errorf("Error executing the compiled binary")
	}
	
	if err := coderunners.SetResourceLimits(runCmd); err != nil {
		return "", fmt.Errorf("Unable to set resource limit: %s", err.Error())
	}

	if _, writeErr := io.WriteString(stdinPipe, stdin); writeErr != nil {
		return "", fmt.Errorf("Error writing to the input pipe")
	}

	stdinPipe.Close()
	if waitErr := runCmd.Wait(); waitErr != nil {
		return "", fmt.Errorf("Resources Limit: Consuming too much resources: %s", waitErr.Error())
	}

	var finalOutput = outputBuffer.String()

	return finalOutput, nil
}
