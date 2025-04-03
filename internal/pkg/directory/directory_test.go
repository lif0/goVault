package directory

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func cleanUp(path string) {
	if parentDir := filepath.Dir(path); parentDir != "." {
		path = parentDir
	}

	os.RemoveAll(path)
}

func TestTryCreateDirsByPath(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		path string

		expectedErr error
	}{
		"path with solo-file": {
			path: "test.log",

			expectedErr: nil,
		},

		"path with file": {
			path: "logs/logs.log",

			expectedErr: nil,
		},

		"path with rw file system with file": {
			path: "/mnt/nonexistent_dir/test.log",

			expectedErr: errors.New("failed to create parent directories for /mnt/nonexistent_dir: mkdir /mnt/nonexistent_dir: permission denied"),
		},

		"path with rw file system": {
			path: "/mnt/nonexistent_dir",

			expectedErr: errors.New("failed to create parent directories for /mnt/nonexistent_dir: mkdir /mnt/nonexistent_dir: permission denied"),
		},

		"path with solo-directory": {
			path: "wal",

			expectedErr: nil,
		},

		"path without file": {
			path: "data/wal",

			expectedErr: nil,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			defer cleanUp(test.path)

			t.Parallel()

			err := TryCreateDirsByPath(test.path)
			assert.Equal(t, test.expectedErr, err)
		})
	}
}

func TestDirectoryExists(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		path     string
		init     func() error
		expected bool
	}{
		"dir is exists": {
			path:     "logs/logs.log",
			init:     func() error { return TryCreateDirsByPath("logs/logs.log") },
			expected: true,
		},

		"dir is not exists": {
			path:     "some_pth/logs.log",
			init:     func() error { return nil },
			expected: false,
		},

		"with error": {
			path: "/mnt/nonexistent_dir/test.log",
			init: func() error {
				_ = TryCreateDirsByPath("/mnt/nonexistent_dir/test.log")
				return nil
			},
			expected: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			defer cleanUp(test.path)

			err := test.init()
			assert.Nil(t, err)
			result := DirectoryExists(test.path)
			assert.Equal(t, test.expected, result)
		})
	}
}
