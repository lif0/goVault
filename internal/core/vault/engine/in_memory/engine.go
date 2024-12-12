package in_memory

import (
	"context"

	"goVault/internal/core/vault/engine"
)

type service struct {
	vault *HashTable
}

func NewEngine() (engine.Engine, error) {
	engine := service{
		vault: NewHashTable(),
	}

	return &engine, nil
}

func (e *service) Set(ctx context.Context, key, value string) {
	e.vault.Set(key, value)
}

func (e *service) Get(ctx context.Context, key string) (string, bool) {
	return e.vault.Get(key)
}

func (e *service) Del(ctx context.Context, key string) {
	e.vault.Del(key)
}
