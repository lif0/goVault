package filesystem

import (
	"fmt"
	"os"

	"goVault/internal/pkg/directory"
)

type SegmentsDirectory struct {
	directory string
}

func NewSegmentsDirectory(directoryPath string) *SegmentsDirectory {
	if !directory.DirectoryExists(directoryPath) {
		directory.TryCreateDirsByPath(directoryPath)
	}

	return &SegmentsDirectory{
		directory: directoryPath,
	}
}

func (d *SegmentsDirectory) ForEach(action func([]byte) error) error {
	files, err := os.ReadDir(d.directory)
	if err != nil {
		return fmt.Errorf("failed to scan directory with segments: %w", err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		filename := fmt.Sprintf("%s/%s", d.directory, file.Name())
		data, err := os.ReadFile(filename)
		if err != nil {
			return err
		}

		if err := action(data); err != nil {
			return err
		}
	}

	return nil
}
