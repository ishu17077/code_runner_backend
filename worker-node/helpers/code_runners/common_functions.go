package coderunners

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"syscall"

	v3 "github.com/containerd/cgroups/v3/cgroup2"
)

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

func SaveFile(fileName string, code string) error {
	err := os.WriteFile(fileName, []byte(code), 0755)
	if err != nil {
		return fmt.Errorf("Cannot save file: %w", err)
	}
	return nil
}

func SetUpCGroup() (*v3.Manager, *os.File) {
	var memeoryLimitBytes int64 = 200 * 1024 * 1024
	var highThresholdBytes int64 = 170 * 1024 * 1024
	var cpuPeriodMicroSec = uint64(1000000)
	var cpuQuotaMicroSec = int64(200000)
	const oomKillEnabledValue = "1"
	// cpuLimitString := fmt.Sprintf("%d %d", cpuQuotaMicroSec, cpuPeriodMicroSec)
	resources := v3.Resources{
		Memory: &v3.Memory{
			High: &highThresholdBytes,
			Max:  &memeoryLimitBytes,
			Swap: &[]int64{0}[0],
		},
		IO: &v3.IO{
			Max: []v3.Entry{
				{
					Major: maj,
					Minor: min,
					Type:  v3.ReadBPS, // Bytes per second
					Rate:  0,
				},
				{
					Major: maj,
					Minor: min,
					Type:  v3.WriteBPS, // Bytes per second
					Rate:  0,
				},
				{
					Major: maj,
					Minor: min,
					Type:  v3.ReadIOPS, // I/O Operations per second
					Rate:  0,
				},
				{
					Major: maj,
					Minor: min,
					Type:  v3.WriteIOPS, // I/O Operations per second
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
	runCmd := exec.Command("rm", "-rf", "/temp/*")
	_, err := runCmd.CombinedOutput()
	if err != nil {
		log.Fatalf("Cannot delete temp directory contents")
	}
}
