package database

import (
	"context"

	"goVault/internal/core/query"
)

type parserLayer interface {
	Transition(query string) (out *query.Query, err error)
}

type vaultLayer interface {
	Set(context.Context, string, string) error
	Get(context.Context, string) (string, error)
	Del(context.Context, string) error
}
