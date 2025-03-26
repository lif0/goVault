package initialization

import (
	"errors"
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"goVault/internal"
	"goVault/internal/configuration"
	"goVault/internal/pkg/dir"
)

const (
	debugLevel = "debug"
	infoLevel  = "info"
	warnLevel  = "warn"
	errorLevel = "error"
)

const (
	defaultEncoding   = "json"
	defaultLevel      = zapcore.InfoLevel
	defaultOutputPath = "goVault.log"
)

func CreateLogger(cfg *configuration.LoggingConfig) (internal.Logger, error) {
	level := defaultLevel
	output := defaultOutputPath

	if cfg != nil {
		if cfg.Level != "" {
			supportedLoggingLevels := map[string]zapcore.Level{
				debugLevel: zapcore.DebugLevel,
				infoLevel:  zapcore.InfoLevel,
				warnLevel:  zapcore.WarnLevel,
				errorLevel: zapcore.ErrorLevel,
			}

			var found bool
			if level, found = supportedLoggingLevels[cfg.Level]; !found {
				return nil, errors.New("logging level is incorrect")
			}
		}

		if cfg.Output != "" {
			err := dir.CreateParentDirIfNeedIt(cfg.Output)
			if err != nil {
				return nil, fmt.Errorf("failed create dir: %s", cfg.Output)
			}

			output = cfg.Output
		}
	}

	loggerCfg := zap.NewProductionConfig()
	loggerCfg.Encoding = defaultEncoding
	loggerCfg.Level = zap.NewAtomicLevelAt(level)
	loggerCfg.OutputPaths = []string{output}
	loggerCfg.EncoderConfig.TimeKey = "time"
	loggerCfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	if cfg != nil && cfg.Stdout {
		loggerCfg.OutputPaths = append(loggerCfg.OutputPaths, "stdout")
	}

	return loggerCfg.Build()
}
