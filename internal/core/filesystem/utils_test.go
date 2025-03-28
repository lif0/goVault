package filesystem

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSegmentCreateFile(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		path        string
		expectedErr error
		cleanUp     func()
	}{
		"success": {
			path:        "0.wal",
			expectedErr: nil,
			cleanUp:     func() { os.Remove("0.wal") },
		},

		"error create file": {
			path:        "root/no_permission.wal",
			expectedErr: errors.New("open root/no_permission.wal: no such file or directory"),
			cleanUp:     func() {},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			defer test.cleanUp()

			t.Parallel()

			file, err := CreateFile(test.path)
			if test.expectedErr != nil {
				require.EqualError(t, err, test.expectedErr.Error())
			} else {
				require.NoError(t, err)
			}

			file.Close()
		})
	}
}
