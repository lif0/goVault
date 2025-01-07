//go:generate mockgen -destination=./../../mocks/net/conn.go -package net_mock net Conn,Listener,Addr
//go:generate mockgen -destination ./../../mocks/network/contract.go -package ${GOPACKAGE}_mock . TCPServer
package network

import (
	"context"
	"net"
	"time"

	"goVault/internal"
	"goVault/internal/concurrency"
)

type TCPHandler = func(context.Context, []byte) []byte
type TCPServerOption func(*server)

type server struct {
	listener  net.Listener
	semaphore concurrency.Semaphore

	idleTimeout    time.Duration
	bufferSize     int
	maxConnections uint

	logger internal.Logger
}

type TCPServer interface {
	HandleQueries(ctx context.Context, handler TCPHandler)
}
