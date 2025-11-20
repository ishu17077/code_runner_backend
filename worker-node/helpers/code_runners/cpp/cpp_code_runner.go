package cpp

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"time"

	coderunners "github.com/ishu17077/code_runner_backend/worker-node/helpers/code_runners"
	"github.com/ishu17077/code_runner_backend/worker-node/models"
)

const filePath = "./temp/main.cpp"
const outputPath = "./temp/main"

func PreCompilationTask(submission models.Submission) error {
	err := coderunners.SaveFile(filePath, submission.Code)
	if err != nil {
		return err
	}

	return nil
}

func CheckSubmission(submission models.Submission, test models.Test) (string, error) {
	res, err := ExecuteCode(outputPath, test.Stdin)
	if err != nil {
		return "FAILED", err
	}

	if strings.TrimSpace(res) == strings.TrimSpace(test.ExpectedOutput) {
		return "SUCCESS", nil
	}

	return "FAILED", fmt.Errorf("FAILED: Expected output: %s. Actual output: %s", test.ExpectedOutput, res)
}

func CompileCode(filePath string, outputPath string) error {
	cmd := exec.Command("g++", filePath, "-o", outputPath)

	_, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Errorf("Compilation  Failed: %s", err.Error())
	}

	return nil
}

func ExecuteCode(binaryFilePath string, stdin string) (string, error) {
	runCmd := exec.Command(binaryFilePath)
	stdinPipe, pipErr := runCmd.StdinPipe()

	if pipErr != nil {
		fmt.Errorf("Error connecting pipe input")
	}

	var outputBuffer bytes.Buffer

	runCmd.Stdin = &outputBuffer

	if startErr := runCmd.Start(); startErr != nil {
		return "", fmt.Errorf("Error starting the program")
	}

	if _, writeErr := io.WriteString(stdinPipe, stdin); writeErr != nil{
		return "", fmt.Errorf("Error writing to the input pipe")
	}

	stdinPipe.Close()
	runCmd.WaitDelay = time.Duration(10*time.Second)

	if waitErr := runCmd.WaitDelay



}
