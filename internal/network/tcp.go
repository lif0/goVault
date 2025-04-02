package network

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"time"

	"go.uber.org/zap"

	"goVault/internal"
	"goVault/internal/pkg/concurrency"
)

func NewTCPServer(address string, l internal.Logger, options ...TCPServerOption) (TCPServer, error) {
	if l == nil {
		return nil, errors.New("logger is invalid")
	}

	listener, err := net.Listen("tcp", address)
	if err != nil {
		return nil, fmt.Errorf("error starting server: %v", err)
	}

	server := &server{
		listener: listener,
		logger:   l,
	}

	for _, option := range options {
		option(server)
	}

	if server.maxConnections != 0 {
		server.semaphore = concurrency.NewSemaphore(server.maxConnections)
	}
	if server.bufferSize == 0 {
		server.bufferSize = 4 << 10 // 4096 bytes
	}

	return server, nil
}

func (s *server) HandleQueries(ctx context.Context, handler TCPHandler) {
	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()

		for {
			conn, err := s.listener.Accept()
			if err != nil {
				if errors.Is(err, net.ErrClosed) {
					return
				}

				s.logger.Error("failed to accept: %w", zap.Error(err))
				continue
			}

			s.semaphore.Acquire()
			go func(connection net.Conn) {
				defer s.semaphore.Release()
				s.handleConnection(ctx, conn, handler)
			}(conn)
		}
	}()

	<-ctx.Done()
	s.listener.Close()

	wg.Wait()
}

func (s *server) handleConnection(ctx context.Context, connection net.Conn, handler TCPHandler) {
	defer func() {
		if v := recover(); v != nil {
			s.logger.Error("captured panic", zap.Any("panic", v))
		}

		if err := connection.Close(); err != nil {
			s.logger.Warn("failed to close connection", zap.Error(err))
		}
	}()

	// reuse buffer for requests
	request := make([]byte, s.bufferSize)

	for {
		if s.idleTimeout != 0 {
			if err := connection.SetReadDeadline(time.Now().Add(s.idleTimeout)); err != nil {
				s.logger.Warn("failed to set read deadline", zap.Error(err))
				break
			}
		}

		count, err := connection.Read(request)
		if err != nil && err != io.EOF {
			s.logger.Warn(
				"failed to read data",
				zap.String("address", connection.RemoteAddr().String()),
				zap.Error(err),
			)
			break
		} else if count == s.bufferSize {
			s.logger.Warn("small buffer size", zap.Int("buffer_size", s.bufferSize))
			break
		}

		if s.idleTimeout != 0 {
			if err := connection.SetWriteDeadline(time.Now().Add(s.idleTimeout)); err != nil {
				s.logger.Warn("failed to set read deadline", zap.Error(err))
				break
			}
		}

		response := handler(ctx, request[:count])
		if _, err := connection.Write(response); err != nil {
			s.logger.Warn(
				"failed to write data",
				zap.String("address", connection.RemoteAddr().String()),
				zap.Error(err),
			)
			break
		}
	}
}
