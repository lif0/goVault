//go:generate mockgen -destination ./../../../mocks/core/vault/contract.go -package ${GOPACKAGE}_mock . Vault
package vault

import "context"

type Vault interface {
	Set(ctx context.Context, key, value string) error
	Get(ctx context.Context, key string) (string, error)
	Del(ctx context.Context, key string) error
}
