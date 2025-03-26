package configuration

import (
	"errors"
	"io"
	"strings"
	"testing"
)

func TestLoad(t *testing.T) {
	tests := []struct {
		name           string
		reader         io.Reader
		expectedError  string
		expectedConfig *Config
	}{
		{
			name:          "Valid configuration",
			reader:        strings.NewReader(validYAML()),
			expectedError: "",
			expectedConfig: &Config{
				Engine: &EngineConfig{Type: "memory"},
				Network: &NetworkConfig{
					Address:        "127.0.0.1:8080",
					MaxConnections: 100,
					MaxMessageSize: "2MB",
					IdleTimeout:    30,
				},
				Logging: &LoggingConfig{
					Level:  "info",
					Output: "stdout",
				},
				WAL: &WALConfig{
					FlushingBatchSize:    100,
					FlushingBatchTimeout: "10ms",
					MaxSegmentSize:       "10MB",
					DataDirectory:        "/goVault/data/wal",
				},
			},
		},
		{
			name:          "Nil reader",
			reader:        nil,
			expectedError: "incorrect reader",
		},
		{
			name:          "Invalid YAML",
			reader:        strings.NewReader("invalid_yaml: : : "),
			expectedError: "failed to parse config",
		},
		{
			name:          "Read error",
			reader:        &errorReader{},
			expectedError: "falied to read buffer",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config, err := Load(tt.reader)

			if tt.expectedError != "" {
				if err == nil || !strings.Contains(err.Error(), tt.expectedError) {
					t.Fatalf("expected error containing %q, got %v", tt.expectedError, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if tt.expectedConfig != nil && !compareConfigs(config, tt.expectedConfig) {
				t.Fatalf("expected config %+v, got %+v", tt.expectedConfig, config)
			}
		})
	}
}

// errorReader is a mock reader that always returns an error.
type errorReader struct{}

func (e *errorReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("falied to read buffer")
}

// validYAML returns a valid YAML string for testing.
func validYAML() string {
	return `
engine:
  type: memory
network:
  address: 127.0.0.1:8080
  max_connections: 100
  max_message_size: 2MB
  idle_timeout: 30
logging:
  level: info
  output: stdout
wal:
  flushing_batch_size: 100
  flushing_batch_timeout: "10ms"
  max_segment_size: "10MB"
  data_directory: "/goVault/data/wal"
`
}

// compareConfigs compares two Config structs for equality.
func compareConfigs(c1, c2 *Config) bool {
	return c1.Engine.Type == c2.Engine.Type &&
		c1.Network.Address == c2.Network.Address &&
		c1.Network.MaxConnections == c2.Network.MaxConnections &&
		c1.Network.MaxMessageSize == c2.Network.MaxMessageSize &&
		c1.Network.IdleTimeout == c2.Network.IdleTimeout &&
		c1.Logging.Level == c2.Logging.Level &&
		c1.Logging.Output == c2.Logging.Output
}
