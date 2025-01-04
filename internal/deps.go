//go:generate mockgen -destination ./../mocks/deps.go -package ${GOPACKAGE}_mock . Logger
package internal

import (
	"go.uber.org/zap"
)

type Logger interface {
	Debug(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
	Warn(msg string, fields ...zap.Field)
	Fatal(msg string, fields ...zap.Field)
}
