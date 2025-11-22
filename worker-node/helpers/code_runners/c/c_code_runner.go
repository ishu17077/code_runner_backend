package c

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

func CheckSubmission(submission models.Submission, test models.TestCase) (enums.CurrentStatus, error) {

	//TODO: Impl executeCcode test case
	res, err := executeCode(outputPath, test.Stdin)
	if err != nil {
		return enums.FAILED, err
	}

	if strings.TrimSpace(res) == strings.TrimSpace(test.ExpectedOutput) {
		return enums.SUCCESS, nil
	}
	return enums.FAILED, fmt.Errorf("FAILED: Expected output: %s. Actual output: %s", test.ExpectedOutput, res)
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

	var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	runCmd := exec.CommandContext(ctx, binaryFilePath)
	// coderunners.SetLimitsAndPermissions(runCmd)
	stdinPipe, pipeErr := runCmd.StdinPipe()
	if pipeErr != nil {
		return "", fmt.Errorf("Error connecting pipe input")
	}
	var outputBuffer bytes.Buffer
	runCmd.Stdout = &outputBuffer

	if startErr := runCmd.Start(); startErr != nil {
		return "", fmt.Errorf("Error starting the program")
	}

	_, writeErr := io.WriteString(stdinPipe, stdin)
	if writeErr != nil {
		return "", fmt.Errorf("Error writing input to stdin")
	}
	stdinPipe.Close()

	if waitErr := runCmd.Wait(); waitErr != nil {
		return "", fmt.Errorf("Error executing file")
	}
	var finalOutput = outputBuffer.String()
	return finalOutput, nil
}
