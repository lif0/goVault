package vault

import (
	"context"
	"errors"
	"fmt"

	"goVault/internal"
)

type vault struct {
	engine Engine
	logger internal.Logger
	wal    WAL
}

func NewVault(e Engine, w WAL, l internal.Logger, options ...VaultOption) (Vault, error) {
	if e == nil {
		return nil, errors.New("engine is invalid")
	}

	if l == nil {
		return nil, errors.New("logger is invalid")
	}

	v := &vault{
		wal:    w,
		engine: e,
		logger: l,
	}

	for _, option := range options {
		option(v)
	}

	return v, nil
}

func (v *vault) Set(ctx context.Context, key, value string) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	if v.wal != nil {
		ch := make(chan struct{})
		v.wal.Write(fmt.Sprintf("SET %s %s", key, value), func() { ch <- struct{}{} })
		<-ch
	}

	v.engine.Set(ctx, key, value)
	return nil
}

func (v *vault) Get(ctx context.Context, key string) (string, error) {
	if ctx.Err() != nil {
		return "", ctx.Err()
	}

	value, found := v.engine.Get(ctx, key)
	if !found {
		return "", ErrVaultNotFound
	}

	return value, nil
}

func (v *vault) Del(ctx context.Context, key string) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	if v.wal != nil {
		ch := make(chan struct{})
		v.wal.Write(fmt.Sprintf("DEL %s", key), func() { ch <- struct{}{} })
		<-ch
	}

	v.engine.Del(ctx, key)
	return nil
}
