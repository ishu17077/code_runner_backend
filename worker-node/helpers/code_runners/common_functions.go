package coderunners

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"

	v2 "github.com/containerd/cgroups/v3/cgroup2"
)

var (
	CGroupFile    *os.File
	CGroupManager *v2.Manager
)

func init() {
	CGroupManager, CGroupFile = setUpCGroup()
}

func SaveFile(fileName string, code string) error {
	err := os.WriteFile(fileName, []byte(code), 0755)
	if err != nil {
		return fmt.Errorf("Cannot save file: %w", err)
	}
	return nil
}

func setUpCGroup() (*v2.Manager, *os.File) {
	var memeoryLimitBytes int64 = 200 * 1024 * 1024
	const cpuPeriodMicroSec = 100000
	const cpuQuotaMicroSec = 80000
	cpuLimitString := fmt.Sprintf("%d %d", cpuQuotaMicroSec, cpuPeriodMicroSec)
	resources := v2.Resources{
		Memory: &v2.Memory{
			Max: &memeoryLimitBytes,
		},
		CPU: &v2.CPU{
			Max: v2.CPUMax(cpuLimitString),
		},
	}

	cGroupPath := filepath.Join("/cgroup", "program_limit")

	manager, err := v2.NewManager(cGroupPath, "", &resources)
	if err != nil {
		log.Fatalf("Error setting new manager for C group")
	}
	cgroupFile, err := os.OpenFile(cGroupPath, os.O_RDONLY, 0)
	if err != nil {
		log.Fatalf("Failed to create a CGroup File")
	}
	return manager, cgroupFile
}

func SetLimitsAndPermissions(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Credential: &syscall.Credential{
			Uid: 6969,
			Gid: 7070,
		},
		UseCgroupFD: true,
		CgroupFD:    int(CGroupFile.Fd()),
	}
}
