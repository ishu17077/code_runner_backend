package coderunners

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	v3 "github.com/containerd/cgroups/v3/cgroup2"
)

var (
	CGroupFile    *os.File
	CGroupManager *v3.Manager
)

// func init() {
// 	CGroupManager, CGroupFile = SetUpCGroup()
// }

func SaveFile(fileName string, code string) error {
	err := os.WriteFile(fileName, []byte(code), 0755)
	if err != nil {
		return fmt.Errorf("Cannot save file: %w", err)
	}
	return nil
}

// func SetUpCGroup() (*v3.Manager, *os.File) {
// 	var memeoryLimitBytes int64 = 200 * 1024 * 1024
// 	var cpuPeriodMicroSec = uint64(100000)
// 	var cpuQuotaMicroSec = int64(80000)
// 	// cpuLimitString := fmt.Sprintf("%d %d", cpuQuotaMicroSec, cpuPeriodMicroSec)
// 	resources := v3.Resources{
// 		Memory: &v3.Memory{
// 			Max: &memeoryLimitBytes,
// 		},
// 		CPU: &v3.CPU{
// 			Max: v3.NewCPUMax(
// 				&cpuQuotaMicroSec,
// 				&cpuPeriodMicroSec,
// 			),
// 		},
// 	}

// 	cGroupPath := "/code_runner"
// 	manager, err := v3.NewManager("/sys/fs/cgroup/", cGroupPath, &resources)
// 	if err != nil {
// 		log.Fatalf("Error setting new manager for C group %s", err.Error())
// 	}
// 	cgroupFile, err := os.OpenFile("/sys/fs/cgroup/code_runner", os.O_RDONLY, 0)
// 	if err != nil {
// 		log.Fatalf("Failed to create a CGroup File")
// 	}
// 	return manager, cgroupFile
// }

// func SetLimitsAndPermissions(cmd *exec.Cmd) {
// 	cmd.SysProcAttr = &syscall.SysProcAttr{
// 		Credential: &syscall.Credential{
// 			Uid: 6969,
// 			Gid: 7070,
// 		},
// 		UseCgroupFD: true,
// 		CgroupFD:    int(CGroupFile.Fd()),
// 	}
// }

func CleanUp() {
	runCmd := exec.Command("rm", "-rf", "/temp/*")
	_, err := runCmd.CombinedOutput()
	if err != nil {
		log.Fatalf("Cannot delete temp directory contents")
	}
}
