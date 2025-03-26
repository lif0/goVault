package wal

import (
	"goVault/internal/configuration"
	"goVault/internal/pkg/unit"
	"sync"
	"time"
)

var (
	defaultFlushingBatchSize    = 100
	defaultFlushingBatchTimeout = int64(200 /*ms*/)
)

type wal struct {
	mu      sync.Mutex
	data    []WALItem
	chFlush chan []WALItem
	walFile *walSegment

	flushingBatchSize     int
	flushingBatchTimeout  int64
	flushingBatchLastTime time.Time
}

func New(cfg configuration.WALConfig) (WAL, error) {
	flushingBatchSize := defaultFlushingBatchSize
	flushingBatchTimeout := defaultFlushingBatchTimeout

	if cfg.FlushingBatchSize == 0 {
		flushingBatchSize = cfg.FlushingBatchSize
	}

	if cfg.FlushingBatchTimeout != "" {
		fbt, err := unit.ParseDuration(cfg.FlushingBatchTimeout)
		if err != nil {
			return nil, err
		}

		flushingBatchTimeout = fbt
	}

	walFile, err := NewWALSegment(cfg)
	if err != nil {
		return nil, err
	}

	service := &wal{
		mu:      sync.Mutex{},
		data:    make([]WALItem, 0, flushingBatchSize),
		chFlush: make(chan []WALItem, 1), // TODO: cover metric size buffer
		walFile: walFile,

		flushingBatchSize:    flushingBatchSize,
		flushingBatchTimeout: flushingBatchTimeout,
	}

	return service, nil
}

func (w *wal) Write(wal string, fCommit func()) error {
	w.mu.Lock()

	if len(w.data) >= w.flushingBatchSize {
		w.chFlush <- w.data
		w.data = make([]WALItem, 0, w.flushingBatchSize)
		w.mu.Unlock()
		return nil
	}

	w.data = append(w.data, WALItem{value: wal, fCommit: fCommit})
	w.mu.Unlock()
	return nil
}

func (w *wal) runFlushHandler() {
	for wals := range w.chFlush {
		for _, v := range wals {
			w.walFile.Write([]byte(v.value))
			v.fCommit()
		}
	}
}
