package initialization

import (
	"errors"

	"goVault/internal"
	"goVault/internal/configuration"
	"goVault/internal/network"
	"goVault/internal/pkg/unit"
)

const defaultServerAddress = ":3231"

func CreateNetwork(cfg *configuration.NetworkConfig, logger internal.Logger) (*network.TCPServer, error) {
	if logger == nil {
		return nil, errors.New("logger is invalid")
	}

	address := defaultServerAddress
	var options []network.TCPServerOption

	if cfg != nil {
		if cfg.Address != "" {
			address = cfg.Address
		}

		if cfg.MaxConnections != 0 {
			options = append(options, network.WithServerMaxConnectionsNumber(uint(cfg.MaxConnections)))
		}

		if cfg.MaxMessageSize != "" {
			size, err := unit.ParseDigitalStorage(cfg.MaxMessageSize)
			if err != nil {
				return nil, errors.New("incorrect max message digital storage")
			}

			options = append(options, network.WithServerBufferSize(uint(size)))
		}

		if cfg.IdleTimeout != 0 {
			options = append(options, network.WithServerIdleTimeout(cfg.IdleTimeout))
		}
	}

	return network.NewTCPServer(address, logger, options...)
}
