package wal

type WAL interface {
	Write(wal string, fCommit func()) error
}

type WALItem struct {
	value   string
	fCommit func()
}
