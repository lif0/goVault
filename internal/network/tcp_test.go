package network

import (
	"context"
	"errors"
	"io"
	"net"
	"testing"
	"time"

	"goVault/internal"
	internal_mock "goVault/mocks"

	"go.uber.org/mock/gomock"
)

func TestNewTCPServer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := internal_mock.NewMockLogger(ctrl)

	tests := []struct {
		name          string
		address       string
		logger        internal.Logger
		options       []TCPServerOption
		expectedError error
	}{
		{
			name:          "Valid server creation",
			address:       "127.0.0.1:1231",
			logger:        logger,
			expectedError: nil,
		},
		{
			name:          "Invalid address",
			address:       "invalid_address",
			logger:        logger,
			expectedError: errors.New("error starting server: listen tcp: address invalid_address: missing port in address"),
		},
		{
			name:          "Nil logger",
			address:       "127.0.0.1:1232",
			logger:        nil,
			expectedError: errors.New("logger is invalid"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server, err := NewTCPServer(tt.address, tt.logger, tt.options...)

			if tt.expectedError != nil {
				if err == nil || err.Error() != tt.expectedError.Error() {
					t.Fatalf("expected error %v, got %v", tt.expectedError, err)
				}
			} else {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
				if server == nil {
					t.Fatalf("expected server instance, got nil")
				}
			}
		})
	}
}

func TestTCPServer_HandleQueries(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := internal_mock.NewMockLogger(ctrl)
	logger.EXPECT().Warn(gomock.Any(), gomock.Any())
	handler := func(ctx context.Context, request []byte) []byte {
		return append([]byte("response: "), request...)
	}

	server, err := NewTCPServer("127.0.0.1:1233", logger)
	if err != nil {
		t.Fatalf("failed to create server: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	errCh := make(chan error, 1)

	go func() {
		conn, err := net.Dial("tcp", server.listener.Addr().String())
		if err != nil {
			errCh <- err
			return
		}
		defer conn.Close()

		_, err = conn.Write([]byte("test request"))
		if err != nil {
			errCh <- err
			return
		}

		response := make([]byte, 1024)
		n, err := conn.Read(response)
		if err != nil && err != io.EOF {
			errCh <- err
			return
		}

		expectedResponse := "response: test request"
		if string(response[:n]) != expectedResponse {
			errCh <- errors.New("unexpected response: " + string(response[:n]))
			return
		}

		errCh <- nil
	}()

	server.HandleQueries(ctx, handler)

	if err := <-errCh; err != nil {
		t.Fatalf("test failed: %v", err)
	}
}
