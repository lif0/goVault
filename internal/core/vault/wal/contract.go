//go:generate mockgen -destination ./../../../../mocks/core/vault/wal/contract.go -package ${GOPACKAGE}_mock . WAL
package wal

import "context"

type WAL interface {
	Start(context.Context)
	Write(wal string, fCommit func()) error
}

type WALItem struct {
	value   string
	fCommit func()
}
