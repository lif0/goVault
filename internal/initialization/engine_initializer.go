package initialization

import (
	"errors"

	"goVault/internal"
	"goVault/internal/configuration"
	"goVault/internal/core/vault/engine"
	"goVault/internal/core/vault/engine/in_memory"
)

func CreateEngine(cfg *configuration.EngineConfig, logger internal.Logger) (engine.Engine, error) {
	if logger == nil {
		return nil, errors.New("logger is invalid")
	}

	if cfg == nil {
		return in_memory.NewEngine(logger)
	}

	if cfg.Type != "" {
		supportedTypes := map[string]struct{}{
			"in_memory": {},
		}

		if _, found := supportedTypes[cfg.Type]; !found {
			return nil, errors.New("engine type is incorrect")
		}
	}

	return in_memory.NewEngine(logger)
}
