package directory

import (
	"fmt"
	"os"
	"path/filepath"
)

// Create parent dirs from path if need it
func TryCreateDirsByPath(path string) error {
	if ext := filepath.Ext(path); ext != "" {
		path = filepath.Dir(path)
	}

	if path != "." {
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			return fmt.Errorf("failed to create parent directories for %s: %v", path, err)
		}
	}

	return nil
}

func DirectoryExists(path string) bool {
	if ext := filepath.Ext(path); ext != "" {
		path = filepath.Dir(path)
	}

	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}

	return err == nil
}
