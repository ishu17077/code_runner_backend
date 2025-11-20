package coderunners

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func SaveFile(fileName string, code string) error {
	err := os.WriteFile(fileName, []byte(code), 0755)
	if err != nil {
		return fmt.Errorf("Cannot save file: %w", err)
	}
	return nil
}

func setLimitsAndPermissions(cmd *exec.Cmd){
	cpuLimit := syscall.Rlimit{
		Cur: 10,
		Max: 10,
	}

	cmd.SysProcAttr = &syscall.SysProcAttr{
		
	}
}