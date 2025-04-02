package initialization

import (
	"errors"
	"fmt"
	"time"

	"goVault/internal"
	"goVault/internal/configuration"
	"goVault/internal/core/filesystem"
	"goVault/internal/core/vault/wal"
	"goVault/internal/pkg/directory"
	"goVault/internal/pkg/unit"
)

const (
	defaultFlushingBatchLength  = 100
	defaultFlushingBatchTimeout = time.Millisecond * 10
	defaultMaxSegmentSize       = 10 << 20 /*10MB*/
	defaultWALDataDirectory     = "./goVault/data/wal"
)

func CreateWAL(cfg *configuration.WALConfig, logger internal.Logger) (wal.WAL, error) {
	if logger == nil {
		return nil, errors.New("logger is invalid")
	} else if cfg == nil {
		return nil, nil
	}

	flushingBatchLength := defaultFlushingBatchLength
	flushingBatchTimeout := defaultFlushingBatchTimeout
	maxSegmentSize := defaultMaxSegmentSize
	dataDirectory := defaultWALDataDirectory

	if cfg.FlushingBatchLength != 0 {
		flushingBatchLength = cfg.FlushingBatchLength
	}

	if cfg.FlushingBatchTimeout != 0 {
		flushingBatchTimeout = cfg.FlushingBatchTimeout
	}

	if cfg.MaxSegmentSize != "" {
		size, err := unit.ParseDigitalStorage(cfg.MaxSegmentSize)
		if err != nil {
			return nil, errors.New("max segment size is incorrect")
		}

		maxSegmentSize = size
	}

	if cfg.DataDirectory != "" {
		dataDirectory = cfg.DataDirectory

		if !directory.DirectoryExists(dataDirectory) {
			if err := directory.TryCreateDirsByPath(dataDirectory); err != nil {
				return nil, fmt.Errorf("fail to create: %v", err)
			}
		}
	}

	segment := filesystem.NewSegment(dataDirectory, maxSegmentSize)
	writer, err := wal.New(segment, flushingBatchLength, flushingBatchTimeout, logger)
	if err != nil {
		return nil, err
	}

	return writer, nil
}
