package vault

import "context"

type Vault interface {
	Set(ctx context.Context, key, value string) error
	Get(ctx context.Context, key string) (string, error)
	Del(ctx context.Context, key string) error

	//ProcessQuery(rawQuery string) Result
}

type Result struct {
	Success bool

	RawValue string
	Type     string
	Error    error
}
