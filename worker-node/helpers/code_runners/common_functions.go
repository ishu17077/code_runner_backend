package coderunners

import (
	"fmt"
	"os"
)

func SaveFile(fileName string, code string) error {
	err := os.WriteFile(fileName, []byte(code), 0755)
	if err != nil {
		return fmt.Errorf("Cannot save file: %w", err)
	}
	return nil
}
