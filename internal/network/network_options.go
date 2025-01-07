package network

import "time"

func WithServerIdleTimeout(timeout time.Duration) TCPServerOption {
	return func(server *server) {
		server.idleTimeout = timeout
	}
}

func WithServerBufferSize(size uint) TCPServerOption {
	return func(server *server) {
		server.bufferSize = int(size)
	}
}

func WithServerMaxConnectionsNumber(count uint) TCPServerOption {
	return func(server *server) {
		server.maxConnections = count
	}
}
