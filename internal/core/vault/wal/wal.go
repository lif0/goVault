package wal

import (
	"context"
	"sync"
	"time"

	"goVault/internal"
	"goVault/internal/core/filesystem"
)

type wal struct {
	logger internal.Logger

	mu      sync.Mutex
	data    []WALItem
	chFlush chan []WALItem
	walFile *filesystem.Segment

	maxBatchSize      int
	flushBatchTimeout time.Duration
}

func New(segment *filesystem.Segment, flushingBatchLength int, flushingBatchTimeout time.Duration, logger internal.Logger) (WAL, error) {
	service := &wal{
		logger: logger,

		mu:      sync.Mutex{},
		data:    make([]WALItem, 0, flushingBatchLength),
		chFlush: make(chan []WALItem, 1), // TODO: cover metric size buffer
		walFile: segment,

		maxBatchSize:      flushingBatchLength,
		flushBatchTimeout: flushingBatchTimeout,
	}

	return service, nil
}

func (w *wal) Write(wal string, fCommit func()) error {
	w.mu.Lock()

	if len(w.data) >= w.maxBatchSize {
		w.chFlush <- w.data
		w.data = make([]WALItem, 0, w.maxBatchSize)
		w.mu.Unlock()
		return nil
	}

	w.data = append(w.data, WALItem{value: wal, fCommit: fCommit})
	w.mu.Unlock()
	return nil
}

func (w *wal) Start(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(w.flushBatchTimeout)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				w.logger.Debug("WAL loop: call ctx.Done()")
				w.flushBatch()
				return
			default:
			}

			select {
			case <-ctx.Done():
				w.logger.Debug("WAL loop: call ctx.Done()")
				w.flushBatch()
				return
			case batch := <-w.chFlush:
				w.logger.Debug("WAL loop: <-w.chFlush")
				for _, v := range batch {
					w.walFile.Write([]byte(v.value))
					v.fCommit()
				}
				ticker.Reset(w.flushBatchTimeout)
			case <-ticker.C:
				w.flushBatch()
			}
		}
	}()
}

func (w *wal) flushBatch() {
	var batch []WALItem

	w.mu.Lock()
	batch = w.data
	w.data = make([]WALItem, 0, w.maxBatchSize)
	w.mu.Unlock()

	if len(batch) != 0 {
		w.chFlush <- batch
	}
}
