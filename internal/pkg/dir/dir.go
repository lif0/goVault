package dir

import (
	"fmt"
	"os"
	"path/filepath"
)

// Create parent dirs from path if need it
func CreateParentDirIfNeedIt(filePath string) error {
	parentDir := filepath.Dir(filePath)

	if parentDir != "." {
		err := os.MkdirAll(parentDir, os.ModePerm)
		if err != nil {
			return fmt.Errorf("failed to create parent directories for %s: %v", filePath, err)
		}
	}

	return nil
}
