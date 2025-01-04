//go:generate mockgen -destination ./../../../../mocks/core/vault/engine/in_memory/contract.go -package ${GOPACKAGE}_mock . Engine
package engine

import "context"

type Engine interface {
	Set(ctx context.Context, key, value string)
	Get(ctx context.Context, key string) (string, bool)
	Del(ctx context.Context, key string)
}
