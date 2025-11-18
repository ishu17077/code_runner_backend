package coderunners

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/ishu17077/code_runner_backend/worker-node/models"
)

func PreCompilation(submission models.Submission) error {
	const fileName = "main.c"
	const directoryName = "./temp/"

	if err := SaveFile(filepath.Join(directoryName, fileName), submission.Code); err != nil {
		return err
	}

	if err := CompileCCode(filepath.Join(directoryName, fileName)); err != nil {
		return err
	}
	return nil
}

func CheckSubmission(submission models.Submission, test models.Test) (string, error) {
	const fileName = "main.c"
	const directoryName = "./temp/"

	//TODO: Impl executeCcode test case
	res, err := ExecuteCCode(filepath.Join(directoryName, "main"), test.Stdin)
	if err != nil {
		return "FAILED", err
	}

	if strings.TrimSpace(res) == strings.TrimSpace(test.ExpectedOutput) {
		return "SUCCESS", nil
	}
	return "FAILED", nil
}

func SaveFile(fileName string, code string) error {
	err := os.WriteFile(fileName, []byte(code), 0755)
	if err != nil {
		return fmt.Errorf("Cannot save file: %w", err)
	}
	return nil
}

func CompileCCode(fileName string) error {
	cmd := exec.Command("gcc", fileName, "-o", "./temp/main")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("Failed to compile C code: %w\nOutput: %s", err, output)
	}
	// fileMode := os.FileMode(0755)
	// if chmodErr := os.Chmod("./temp/main", fileMode); chmodErr != nil {
	// 	return fmt.Errorf("Failed to set execute permissions to file")
	// }
	return nil
}

func ExecuteCCode(binaryFileName string, stdin string) (string, error) {
	runCmd := exec.Command("./" + binaryFileName)
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
