package coderunners

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"syscall"

	v3 "github.com/containerd/cgroups/v3/cgroup2"
)

//TODO: Impl compiletime security too

var (
	CGroupFile    *os.File
	CGroupManager *v3.Manager
)

const (
	maj int64 = 8
	min int64 = 0
)

func init() {
	CGroupManager, CGroupFile = SetUpCGroup()
}

func SaveFile(filePath string, code string) error {
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
	if err := SetResourceLimits(runCmd); err != nil {
		return "", fmt.Errorf("Unable to set resource limit: %s", err.Error())
	}

	if _, err := io.WriteString(stdinPipe, stdin); err != nil {
		return "", fmt.Errorf("Error writing to stdin: %s", err.Error())
	}
	stdinPipe.Close()

	if waitErr := runCmd.Wait(); waitErr != nil {
		//? If the command context timed out, rtime limit exceeded.
		if errors.Is(waitErr, context.DeadlineExceeded) {
			return "", fmt.Errorf("Time Limit Exceeded")
		}
		//? If process exited with an ExitError, inspect the  wait status to detect signals.
		if exitErr, ok := waitErr.(*exec.ExitError); ok {
			if status, ok := exitErr.Sys().(syscall.WaitStatus); ok {
				if status.Signaled() {
					if status.Signal() == syscall.SIGKILL {
						return "", fmt.Errorf("Time Limit Exceeded")
					}
					return "", fmt.Errorf("Process killed by signal: %s", status.Signal())
				}
				return "", fmt.Errorf("Process exited with code %d", status.ExitStatus())
			}
		}
		// Fallback message for other errors.
		return "", fmt.Errorf("Resources Limit: Consuming too much resources: %s", waitErr.Error())
	}

	return outputBuffer.String(), nil
}

func SetUpCGroup() (*v3.Manager, *os.File) {
	var memeoryLimitBytes int64 = 200 * 1024 * 1024
	var cpuPeriodMicroSec = uint64(1000000)
	var cpuQuotaMicroSec = int64(200000)
	const oomKillEnabledValue = "1"
	// cpuLimitString := fmt.Sprintf("%d %d", cpuQuotaMicroSec, cpuPeriodMicroSec)
	resources := v3.Resources{
		Memory: &v3.Memory{
			Max:  &memeoryLimitBytes,
			Swap: &[]int64{0}[0],
		},
		IO: &v3.IO{
			Max: []v3.Entry{
				{ //? 0 ops per second in all bps(bytes oer sec) and iops
					Major: maj,
					Minor: min,
					Type:  v3.ReadBPS,
					Rate:  0,
				},
				{
					Major: maj,
					Minor: min,
					Type:  v3.WriteBPS,
					Rate:  0,
				},
				{
					Major: maj,
					Minor: min,
					Type:  v3.ReadIOPS,
					Rate:  0,
				},
				{
					Major: maj,
					Minor: min,
					Type:  v3.WriteIOPS,
					Rate:  0,
				},
			},
		},

		CPU: &v3.CPU{

			Max: v3.NewCPUMax(
				&cpuQuotaMicroSec,
				&cpuPeriodMicroSec,
			),
		},
	}

	cGroupPath := "/code_runner"
	manager, err := v3.NewManager("/cgroup/", cGroupPath, &resources)
	if err != nil {
		log.Fatalf("Error setting new manager for C group %s", err.Error())
	}
	cgroupFile, err := os.OpenFile("/cgroup/code_runner", os.O_RDONLY, 0)
	if err != nil {
		log.Fatalf("Failed to create a CGroup File")
	}
	if err := os.WriteFile("/cgroup/code_runner/memory.oom.group", []byte(oomKillEnabledValue), 0644); err != nil {
		log.Fatalf("Error setting up OOM killer")
	}
	return manager, cgroupFile
}

func SetPermissions(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Credential: &syscall.Credential{
			Uid: 6969,
			Gid: 7070,
		},
	}
	// CGroupManager.AddProc(syscall.Process)
}

func SetResourceLimits(cmd *exec.Cmd) error {

	if err := CGroupManager.AddProc(uint64(cmd.Process.Pid)); err != nil {
		fmt.Printf("Error adding process to cgroup: %v\n", err)
		cmd.Process.Kill()
		return fmt.Errorf("Unable to attach to cgroup")
	}
	return nil
}

func CleanUp() {
	runCmd := exec.Command("sh", "-c", "rm -rf /temp/*")
	_, err := runCmd.CombinedOutput()
	if err != nil {
		log.Fatalf("Cannot delete temp directory contents")
	}
}
