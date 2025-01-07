package initialization

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"goVault/internal/configuration"
)

func cleanUp(cfg *configuration.LoggingConfig) {
	if cfg == nil {
		parentDir := filepath.Dir(defaultOutputPath)

		if parentDir != "." {

		} else {
			_ = os.Remove(defaultOutputPath)
		}
	} else {
		parentDir := filepath.Dir(cfg.Output)

		if parentDir != "." {
			_ = os.RemoveAll(parentDir)
		} else {
			_ = os.Remove(cfg.Output)
		}
	}
}

func TestCreateLogger(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		cfg            *configuration.LoggingConfig
		expectedErr    error
		expectedNilObj bool
	}{
		"create logger without config": {
			cfg:         nil,
			expectedErr: nil,
		},
		"create logger with empty config fields": {
			cfg:         &configuration.LoggingConfig{},
			expectedErr: nil,
		},
		"create logger with config fields": {
			cfg: &configuration.LoggingConfig{
				Level:  debugLevel,
				Output: "test.log",
				Stdout: true,
			},
			expectedErr: nil,
		},
		"create logger with full output path": {
			cfg: &configuration.LoggingConfig{
				Level:  debugLevel,
				Output: "logs/logs.log",
				Stdout: true,
			},
			expectedErr: nil,
		},
		"create logger with broken output path": {
			cfg: &configuration.LoggingConfig{
				Level:  debugLevel,
				Output: "/mnt/nonexistent_dir/test.log",
			},
			expectedErr:    errors.New("failed create dir: /mnt/nonexistent_dir/test.log"),
			expectedNilObj: true,
		},
		"create logger with incorrect level": {
			cfg:            &configuration.LoggingConfig{Level: "incorrect"},
			expectedErr:    errors.New("logging level is incorrect"),
			expectedNilObj: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			defer cleanUp(test.cfg)

			t.Parallel()

			logger, err := CreateLogger(test.cfg)
			assert.Equal(t, test.expectedErr, err)
			if test.expectedNilObj {
				assert.Nil(t, logger)
			} else {
				assert.NotNil(t, logger)
			}
		})
	}
}
