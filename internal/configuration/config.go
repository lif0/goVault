package configuration

import (
	"errors"
	"fmt"
	"io"
	"time"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Engine  *EngineConfig  `yaml:"engine"`
	Network *NetworkConfig `yaml:"network"`
	Logging *LoggingConfig `yaml:"logging"`
}

type EngineConfig struct {
	Type string `yaml:"type"`
}

type NetworkConfig struct {
	Address        string        `yaml:"address"`
	MaxConnections int           `yaml:"max_connections"`
	MaxMessageSize string        `yaml:"max_message_size"`
	IdleTimeout    time.Duration `yaml:"idle_timeout"`
}

type LoggingConfig struct {
	Level  string `yaml:"level"`
	Output string `yaml:"output"`
	Stdout bool   `yaml:"stdout"`
}

func Load(reader io.Reader) (*Config, error) {
	if reader == nil {
		return nil, errors.New("incorrect reader")
	}

	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, errors.New("falied to read buffer")
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return &config, nil
}
