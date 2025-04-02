package filesystem

import (
	"os"
)

func CreateFile(filename string) (*os.File, error) {
	flags := os.O_CREATE | os.O_WRONLY
	file, err := os.OpenFile(filename, flags, 0644) //(rw-, r--, r--)
	if err != nil {
		return nil, err
	}

	return file, err
}

func WriteFile(file *os.File, data []byte, newLine bool) (int, error) {
	writtenBytes, err := file.Write(data)
	if err != nil {
		return 0, err
	}

	if newLine {
		_, _ = file.Write([]byte("\n")) // TODO: a potential location for a bug
	}

	if err = file.Sync(); err != nil {
		return 0, err
	}

	return writtenBytes, nil
}
