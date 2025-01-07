package initialization

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"goVault/internal/configuration"
	internal_mock "goVault/mocks"
	engine_mock "goVault/mocks/core/vault/engine/in_memory"
	network_mock "goVault/mocks/network"
)

func TestInitializerNew(t *testing.T) {

	tests := map[string]struct {
		cfg            *configuration.Config
		expectedErr    error
		expectedNilObj bool
	}{
		"successs": {
			cfg: &configuration.Config{
				Network: &configuration.NetworkConfig{
					Address: "localhost:6666",
				},
			},
			expectedErr:    nil,
			expectedNilObj: false,
		},
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
			if err != nil && test.expectedErr != nil {
				assert.Equal(t, test.expectedErr.Error(), err.Error())
			} else if err != nil || test.expectedErr != nil {
				assert.Equal(t, test.expectedErr, err)
			}

			if test.expectedNilObj {
				assert.Nil(t, initializer)
			} else {
				assert.NotNil(t, initializer)
			}

			cleanUp(nil) // remove log-file
		})
	}
}

func TestInitializerStartDatabase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name        string
		setupMocks  func() *Initializer
		expectedErr error
	}{
		{
			name: "successful",
			setupMocks: func() *Initializer {
				serverMock := network_mock.NewMockTCPServer(ctrl)
				serverMock.EXPECT().HandleQueries(gomock.Any(), gomock.Any())

				i := Initializer{
					logger: internal_mock.NewMockLogger(ctrl),
					engine: engine_mock.NewMockEngine(ctrl),
					server: nil,
				}

				i.server = serverMock
				return &i
			},
			expectedErr: nil,
		},
		{
			name: "query.NewParser return error",
			setupMocks: func() *Initializer {
				i := Initializer{
					logger: nil,
					engine: engine_mock.NewMockEngine(ctrl),
					server: network_mock.NewMockTCPServer(ctrl),
				}

				return &i
			},
			expectedErr: errors.New("logger is invalid"),
		},
		{
			name: "vault.NewVault return error",
			setupMocks: func() *Initializer {
				i := Initializer{
					logger: internal_mock.NewMockLogger(ctrl),
					engine: nil,
				}

				return &i
			},
			expectedErr: errors.New("engine is invalid"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()

			init := tt.setupMocks()

			err := init.StartDatabase(ctx)
			if err != nil && tt.expectedErr != nil {
				assert.Equal(t, tt.expectedErr.Error(), err.Error())
			} else if err != nil || tt.expectedErr != nil {
				assert.Equal(t, tt.expectedErr, err)
			}
		})
	}
}
