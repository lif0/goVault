package in_memory

import (
	"context"

	"goVault/internal"
	"goVault/internal/core/vault/engine"
)

type service struct {
	vault  *HashTable
	logger internal.Logger
}

func NewEngine(logger internal.Logger) (engine.Engine, error) {
	engine := service{
		vault:  NewHashTable(),
		logger: logger,
	}

	return &engine, nil
}

func (e *service) Set(ctx context.Context, key, value string) {
	e.vault.Set(key, value)
	e.logger.Debug("successfull set query")
}

func (e *service) Get(ctx context.Context, key string) (string, bool) {
	v, ok := e.vault.Get(key)
	e.logger.Debug("successfull get query")
	return v, ok
}

func (e *service) Del(ctx context.Context, key string) {
	e.vault.Del(key)
	e.logger.Debug("successfull del query")
}
