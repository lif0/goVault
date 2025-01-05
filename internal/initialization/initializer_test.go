package initialization

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"goVault/internal/configuration"
)

func TestInitializer(t *testing.T) {
	t.Parallel()

	initializer, err := NewInitializer(&configuration.Config{
		Network: &configuration.NetworkConfig{
			Address: "localhost:6666",
		},
	})
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*100)
	defer cancel()

	err = initializer.StartDatabase(ctx)
	require.NoError(t, err)
}

func TestInitializerCases(t *testing.T) {

	tests := map[string]struct {
		cfg            *configuration.Config
		expectedErr    error
		expectedNilObj bool
	}{
		"nil cfg": {
			cfg:            nil,
			expectedErr:    errors.New("failed to initialize: config is invalid"),
			expectedNilObj: true,
		},
		"broken logger": {
			cfg: &configuration.Config{
				Logging: &configuration.LoggingConfig{
					Level: "-100",
				},
			},
			expectedErr:    errors.New("failed to initialize logger: logging level is incorrect"),
			expectedNilObj: true,
		},
		"logger with stdout": {
			cfg: &configuration.Config{
				Logging: &configuration.LoggingConfig{
					Level:  "-100",
					Stdout: false,
				},
			},
			expectedErr:    errors.New("failed to initialize logger: logging level is incorrect"),
			expectedNilObj: true,
		},
		"broken engine": {
			cfg: &configuration.Config{
				Engine: &configuration.EngineConfig{
					Type: "broken_engine_name",
				},
			},
			expectedErr:    errors.New("failed to initialize engine: engine type is incorrect"),
			expectedNilObj: true,
		},
		"broken network": {
			cfg: &configuration.Config{
				Network: &configuration.NetworkConfig{
					MaxMessageSize: "13gg",
				},
			},
			expectedErr:    errors.New("failed to initialize network: incorrect max message digital storage"),
			expectedNilObj: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			initializer, err := NewInitializer(test.cfg)
			assert.Equal(t, test.expectedErr.Error(), err.Error())

			if test.expectedNilObj {
				assert.Nil(t, initializer)
			} else {
				assert.NotNil(t, initializer)
			}
		})
	}
}
