package coderunners

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"syscall"

	currentstatus "github.com/ishu17077/code_runner_backend/worker-node/models/enums/current_status"
)

var TleError error = fmt.Errorf("Time Limit Exceeded")

// var (
// 	cGroupFile    *os.File
// 	cGroupManager *v3.Manager
// )

const (
	maj int64 = 8
	min int64 = 0
)

// func init() {
// 	cGroupManager, cGroupFile = setUpCGroup()
// }

func SaveFile(filePath string, dirPath string, code string) error {
	//! You know it i know it
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		return fmt.Errorf("Cannot create new directory: %w", err)
	}

	if err := os.Chown(dirPath, 6969, 7070); err != nil {
		return fmt.Errorf("Error chowining directory: %w", err)
	}

	err := os.WriteFile(filePath, []byte(code), 0755)
	if err != nil {
		return fmt.Errorf("Cannot save file: %w", err)
	}
	return nil
}

func RunCommandWithInput(runCmd *exec.Cmd, stdin string) (string, error) {
	SetPermissions(runCmd)
	stdinPipe, pipeErr := runCmd.StdinPipe()
	if pipeErr != nil {
		return "", fmt.Errorf("Error connecting pipe input")
	}

	var outputBuffer bytes.Buffer
	runCmd.Stdout = &outputBuffer

	if startErr := runCmd.Start(); startErr != nil {
		return "", fmt.Errorf("Unable to start the program %s", startErr.Error())
	}

	if _, err := io.WriteString(stdinPipe, stdin); err != nil {
		return "", fmt.Errorf("Error writing to stdin: %s", err.Error())
	}
	stdinPipe.Close()

	if waitErr := runCmd.Wait(); waitErr != nil {
		//? If the command context timed out, rtime limit exceeded.
		if errors.Is(waitErr, context.DeadlineExceeded) {

			return "", TleError
		}
		//? If process exited with an ExitError, inspect the  wait status to detect signals.
		if exitErr, ok := waitErr.(*exec.ExitError); ok {
			if status, ok := exitErr.Sys().(syscall.WaitStatus); ok {
				if status.Signaled() {
					if status.Signal() == syscall.SIGKILL {
						return "", TleError
					}
					return "", fmt.Errorf("Process killed by signal: %s", status.Signal())
				}
				return "", fmt.Errorf("Process exited with code %d", status.ExitStatus())
			}
		}
		return "", fmt.Errorf("Resources Limit: Consuming too much resources: %s", waitErr.Error())
	}

	return outputBuffer.String(), nil
}

// func setUpCGroup() (*v3.Manager, *os.File) {
// 	var memeoryLimitBytes int64 = 180 * 1024 * 1024
// 	var highThresholdBytes int64 = 120 * 1024 * 1024
// 	var cpuPeriodMicroSec = uint64(1000000)
// 	var cpuQuotaMicroSec = int64(800000)
// 	const oomKillEnabledValue = "1"
// 	// cpuLimitString := fmt.Sprintf("%d %d", cpuQuotaMicroSec, cpuPeriodMicroSec)
// 	resources := v3.Resources{
// 		Memory: &v3.Memory{
// 			High: &highThresholdBytes,
// 			Max:  &memeoryLimitBytes,
// 			Swap: &[]int64{0}[0],
// 		},
// 		IO: &v3.IO{
// 			Max: []v3.Entry{
// 				{
// 					Major: maj,
// 					Minor: min,
// 					Type:  v3.ReadBPS, // Bytes per second
// 					Rate:  0,
// 				},
// 				{
// 					Major: maj,
// 					Minor: min,
// 					Type:  v3.WriteBPS, // Bytes per second
// 					Rate:  0,
// 				},
// 				{
// 					Major: maj,
// 					Minor: min,
// 					Type:  v3.ReadIOPS, // I/O Operations per second
// 					Rate:  0,
// 				},
// 				{
// 					Major: maj,
// 					Minor: min,
// 					Type:  v3.WriteIOPS, // I/O Operations per second
// 					Rate:  0,
// 				},
// 			},
// 		},

// 		CPU: &v3.CPU{

// 			Max: v3.NewCPUMax(
// 				&cpuQuotaMicroSec,
// 				&cpuPeriodMicroSec,
// 			),
// 		},
// 	}

// 	cGroupPath := "/code_runner"
// 	manager, err := v3.NewManager("/cgroup/", cGroupPath, &resources)
// 	if err != nil {
// 		log.Fatalf("Error setting new manager for C group %s", err.Error())
// 	}
// 	cgroupFile, err := os.OpenFile("/cgroup/code_runner", os.O_RDONLY, 0)
// 	if err != nil {
// 		log.Fatalf("Failed to create a CGroup File")
// 	}
// 	if err := os.WriteFile("/cgroup/code_runner/memory.oom.group", []byte(oomKillEnabledValue), 0644); err != nil {
// 		log.Fatalf("Error setting up OOM killer")
// 	}
// 	return manager, cgroupFile
// }

func SetPermissions(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Credential: &syscall.Credential{
			Uid: 6969,
			Gid: 7070,
		},
	}
	// CGroupManager.AddProc(syscall.Process)
}

// func SetResourceLimits(cmd *exec.Cmd) error {

// 	if err := cGroupManager.AddProc(uint64(cmd.Process.Pid)); err != nil {
// 		fmt.Printf("Error adding process to cgroup: %v\n", err)
// 		cmd.Process.Kill()
// 		return fmt.Errorf("Unable to attach to cgroup")
// 	}
// 	return nil
// }

func CheckOutput(actualOutput string, expectedOutput string) (currentstatus.CurrentStatus, error) {
	actualOutput = strings.TrimSpace(actualOutput)
	expectedOutput = strings.TrimSpace(expectedOutput)

	expectedLines := strings.Split(expectedOutput, "\n")

	outputLines := []string{}
	if actualOutput != "" {
		outputLines = strings.Split(actualOutput, "\n")
	}
	//! Impossible test pass scenario
	if len(outputLines) < len(expectedLines) {
		return currentstatus.FAILED, fmt.Errorf("FAILED: Expected output: %s. Actual output: %s", expectedOutput, actualOutput)
	}

	start := len(outputLines) - len(expectedLines)
	for i, expectedLine := range expectedLines {
		if strings.TrimSpace(expectedLine) != strings.TrimSpace(outputLines[start+i]) {
			return currentstatus.FAILED, fmt.Errorf("FAILED: Expected output: %s. Actual output: %s", expectedOutput, actualOutput)
		}
	}
	return currentstatus.SUCCESS, nil
}

func CleanUp(path string) {
	actualCmd := fmt.Sprintf("rm -rf %s", path)
	runCmd := exec.Command("sh", "-c", actualCmd)
	_, err := runCmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Cannot delete temp directory contents")
	}
}
