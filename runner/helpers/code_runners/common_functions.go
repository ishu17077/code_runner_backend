package coderunners

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"syscall"

	"github.com/ishu17077/code_runner_backend/models"
	currentstatus "github.com/ishu17077/code_runner_backend/models/enums/current_status"
)

var TleError error = fmt.Errorf("Time Limit Exceeded")

// var (
// 	cGroupFile    *os.File
// 	cGroupManager *v3.Manager
// )

// const (
// 	maj int64 = 8
// 	min int64 = 0
// )

// func init() {
// 	cGroupManager, cGroupFile = setUpCGroup()
// }

func SaveFile(filePath string, dirPath string, code string) error {
	//! You know it i know it
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		return fmt.Errorf("Cannot create new directory: %w", err)
	}
	if err := os.Chown(dirPath, 6969, 7070); err != nil {
		return fmt.Errorf("Error chowning file: %w", err)
	}
	err := os.WriteFile(filePath, []byte(code), 0755)
	if err != nil {
		return fmt.Errorf("Cannot save file: %w", err)
	}
	return nil
}

func RunCommandWithInput(runCmd *exec.Cmd, stdin string) (string, error) {
	if runCmd.Args[0] != "-XX:+UseSerialGC" {
		SetPermissions(runCmd)
	}

	// SetPermissions(runCmd)

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
		//? If the command context timed out, time limit exceeded.
		return "", TleError
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
	cmd.Env = []string{
		"PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/root/.dotnet/tools:/opt/Rust/.cargo/bin",
		"HOME=/tmp",
		"RUST_HOME=/opt/Rust",
		"RUSTUP_HOME=/opt/Rust/.rustup",
		"CARGO_HOME=/opt/Rust/.cargo",
		"REALLY=GOOD_LUCK_GETTING_ANYTHING"}
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Credential: &syscall.Credential{
			Uid: 6969,
			Gid: 7070,
		},
	}
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

func CheckJavaOutput(allOutput string, allTests []models.TestCase) (bool, models.Result, error) {
	var result models.Result
	var allPassed = true

	result = extractJsonFromStdout(allOutput)

	if len(result.Results) > len(allTests) {
		return false, emptyResult(currentstatus.INTERNAL_ERROR, "Invalid Results"), fmt.Errorf("Invalid Results")
	}
	if result.Error != "" || result.Status == currentstatus.INTERNAL_ERROR.ToString() {
		return false, result, fmt.Errorf("%s", result.Error)
	}

	for i, execResult := range result.Results {
		if strings.TrimSpace(allTests[i].ExpectedOutput) != strings.TrimSpace(execResult.Status.Stdout) {
			allPassed = false
			result.Status = currentstatus.FAILED.ToString()
			result.Results[i].Status.Message = fmt.Sprintf("Expected Output: %s, Actual Ouptut: %s", allTests[i].ExpectedOutput, execResult.Status.Stdout)
		} else {
			result.Results[i].Status.Current_status = currentstatus.SUCCESS.ToString()
		}
	}
	return allPassed, result, nil
}

func extractJsonFromStdout(res string) models.Result {

	var regExpMatch = regexp.MustCompile(`(?s)---JSON_START---(.*?)---JSON_END---`)
	matches := regExpMatch.FindStringSubmatch(res)
	//? matches[0] is entire block and matches[1] is is content in b/w start and end
	if len(matches) < 2 {
		return emptyResult(currentstatus.RESOURCE_LIMIT_EXCEEDED, "Error consuming too much resources")
	}
	res = strings.TrimSpace(matches[1])

	var result models.Result
	jsonData, err := base64.StdEncoding.DecodeString(res)
	if err != nil {
		return emptyResult(currentstatus.INTERNAL_ERROR, "Unable to execute the program.")

	}
	if err = json.Unmarshal(jsonData, &result); err != nil {
		return emptyResult(currentstatus.INTERNAL_ERROR, "Unable to execute the program.")
	}

	return result

}

func emptyResult(status currentstatus.CurrentStatus, err string) models.Result {
	return models.Result{
		Status:  status.ToString(),
		Results: []models.ExecResult{},
		Error:   err,
	}
}

func CleanUp(path string) {
	actualCmd := fmt.Sprintf("rm -rf %s", path)
	runCmd := exec.Command("sh", "-c", actualCmd)
	_, err := runCmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Cannot delete temp directory contents")
	}
}
