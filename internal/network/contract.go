package network

import (
	"context"
	"net"
	"time"

	"goVault/internal"
	"goVault/internal/concurrency"
)

type TCPHandler = func(context.Context, []byte) []byte
type TCPServerOption func(*TCPServer)

type TCPServer struct {
	listener  net.Listener
	semaphore concurrency.Semaphore

	idleTimeout    time.Duration
	bufferSize     int
	maxConnections int

	logger internal.Logger
}
