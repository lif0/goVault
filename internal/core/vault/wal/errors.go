package wal

import "errors"

var (
	ErrWalStorageErrorWriteDataToDisk = errors.New("wal storage: error write data to disk")
	ErrWalStorageErrorSyncFile        = errors.New("wal storage: file sync error")
)
