package vault

import (
	"context"

	"goVault/internal/core/query"
)

type Parser interface {
	Transition(in string) (out *query.Query, err error)
}

type Engine interface {
	Set(ctx context.Context, key, value string)
	Get(ctx context.Context, key string) (string, bool)
	Del(ctx context.Context, key string)
}
