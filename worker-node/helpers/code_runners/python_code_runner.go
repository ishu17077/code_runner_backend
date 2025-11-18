package coderunners

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/ishu17077/code_runner_backend/worker-node/models"
)

func CheckPythonSubmission(submission models.Submission, test models.Test) (string, error) {
	const fileName = "main.py"
	const directoryName = "./test"

	fileWithDir := filepath.Join(directoryName, fileName)

	if err := SaveFile(fileWithDir, submission.Code); err != nil {
		return "", fmt.Errorf("Error saving the file: %s", err.Error())
	}

	res, err := ExecutePythonCode(fileWithDir, test.Stdin)
	if err != nil {
		return "FAILED", fmt.Errorf("The test was unsuccessful: %s", err.Error())
	}

	if strings.TrimSpace(res) == strings.TrimSpace(test.ExpectedOutput) {
		return "SUCCESS", nil
	}
	return "FAILED", fmt.Errorf("Test: #%s Failed", test.Test_id)
}

func ExecutePythonCode(filePath string, stdin string) (string, error) {
	runCmd := exec.Command("./" + filePath)
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
