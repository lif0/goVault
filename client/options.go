package client

import "time"

const defaultBufferSize = 4 << 10

type TCPClientOption func(*TCPClient)

func WithClientIdleTimeout(timeout time.Duration) TCPClientOption {
	return func(client *TCPClient) {
		client.idleTimeout = timeout
	}
}

func WithClientBufferSize(size uint) TCPClientOption {
	return func(client *TCPClient) {
		client.bufferSize = int(size)
	}
}
