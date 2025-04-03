package filesystem

import (
	"os"
)

var newLineByte = []byte("\n")

// CreateFile opens the file specified by filename in write-only mode with permissions set to 0644.
// If the file does not exist, it will be created. It returns the opened file and any error encountered.
func CreateFile(filename string) (*os.File, error) {
	flags := os.O_CREATE | os.O_WRONLY
	file, err := os.OpenFile(filename, flags, 0644) //(rw-, r--, r--)
	if err != nil {
		return nil, err
	}

	return file, err
}

// WriteFile writes the given data to the specified file. If newLine is true, it writes a newline
// before writing the data. The function ensures all data is synced to storage.
// It returns the number of bytes written and any error encountered.
func WriteFile(file *os.File, data []byte, newLine bool) (int, error) {
	if newLine {
		if _, err := file.Write(newLineByte); err != nil {
			return 0, err
		}
	}

	writtenBytes, err := file.Write(data)
	if err != nil {
		return 0, err
	}

	if err = file.Sync(); err != nil {
		return 0, err
	}

	return writtenBytes, nil
}
