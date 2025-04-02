package network

import (
	"context"
	"errors"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"

	"goVault/internal"
	"goVault/internal/pkg/concurrency"
	internal_mock "goVault/mocks"
	net_mock "goVault/mocks/net"
)

func TestNewTCPServer(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := internal_mock.NewMockLogger(ctrl)

	tests := []struct {
		name        string
		address     string
		logger      internal.Logger
		options     []TCPServerOption
		expectedErr error
	}{
		{
			name:        "Successful",
			address:     "127.0.0.1:0",
			logger:      mockLogger,
			options:     []TCPServerOption{WithServerMaxConnectionsNumber(uint(1))},
			expectedErr: nil,
		},
		{
			name:        "Invalid logger",
			address:     "127.0.0.1:0",
			logger:      nil,
			expectedErr: errors.New("logger is invalid"),
		},
		{
			name:        "Invalid address",
			address:     "invalid_address",
			logger:      mockLogger,
			expectedErr: errors.New("error starting server"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			server, err := NewTCPServer(test.address, test.logger, test.options...)
			if test.expectedErr != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), test.expectedErr.Error())
				assert.Nil(t, server)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, server)
			}
		})
	}
}

func Test_HandleQueries(t *testing.T) {
	ctrl := gomock.NewController(t)

	tests := []struct {
		name       string
		setupMocks func() (*server, net.Conn)
	}{
		{
			name: "Successful query handling",
			setupMocks: func() (*server, net.Conn) {
				mockLogger := internal_mock.NewMockLogger(ctrl)
				mockListener := net_mock.NewMockListener(ctrl)
				mockConn := net_mock.NewMockConn(ctrl)

				s := &server{
					logger:     mockLogger,
					listener:   mockListener,
					bufferSize: 1024,
					semaphore:  concurrency.NewSemaphore(2),
				}

				// Setup mocks for listener
				mockListener.EXPECT().Accept().Return(mockConn, nil).Times(1)
				mockListener.EXPECT().Accept().Return(nil, net.ErrClosed).Times(1)
				mockListener.EXPECT().Close().Return(nil).Times(1)

				// Setup mocks for net.Addr
				mockAddr := net_mock.NewMockAddr(ctrl)
				mockAddr.EXPECT().String().Return("127.0.0.1:12345").Times(1)

				// Setup mocks for connection
				mockConn.EXPECT().Read(gomock.Any()).DoAndReturn(func(b []byte) (int, error) {
					copy(b, []byte("request"))
					return len("request"), nil
				}).Times(1)
				mockConn.EXPECT().RemoteAddr().Return(mockAddr).Times(1)
				mockConn.EXPECT().Write([]byte("response")).Return(len("response"), nil).Times(1)
				mockConn.EXPECT().Close().Return(nil).Times(1)

				// Setup mocks for finish
				// second read will be end with error because i want break infinity cycle from handleConnection
				mockConn.EXPECT().Read(gomock.Any()).Return(0, errors.New("error for exit")).Times(1)
				mockLogger.EXPECT().Warn(gomock.Any(), gomock.Any()).Times(1)

				return s, mockConn
			},
		},
		{
			name: "listener.Accept returns some error",
			setupMocks: func() (*server, net.Conn) {
				mockLogger := internal_mock.NewMockLogger(ctrl)
				mockListener := net_mock.NewMockListener(ctrl)
				mockConn := net_mock.NewMockConn(ctrl)

				s := &server{
					logger:     mockLogger,
					listener:   mockListener,
					bufferSize: 1024,
					semaphore:  concurrency.NewSemaphore(2),
				}

				// Setup mocks for listener
				mockListener.EXPECT().Accept().Return(nil, errors.New("some error")).Times(1)
				mockListener.EXPECT().Accept().Return(nil, net.ErrClosed).Times(1)
				mockListener.EXPECT().Close().Return(nil).Times(1)

				// Setup mocks for logger
				mockLogger.EXPECT().Error(gomock.Any(), gomock.Any()).Times(1)

				return s, mockConn
			},
		},
		{
			name: "listener.Accept returns net.ErrClosed error",
			setupMocks: func() (*server, net.Conn) {
				mockLogger := internal_mock.NewMockLogger(ctrl)
				mockListener := net_mock.NewMockListener(ctrl)
				mockConn := net_mock.NewMockConn(ctrl)

				s := &server{
					logger:     mockLogger,
					listener:   mockListener,
					bufferSize: 1024,
					semaphore:  concurrency.NewSemaphore(1),
				}

				// Setup mocks for listener
				mockListener.EXPECT().Accept().Return(nil, net.ErrClosed).Times(1)
				mockListener.EXPECT().Close().Return(nil).Times(1)

				return s, mockConn
			},
		},
		{
			name: "conn.SetReadDeadline return error",
			setupMocks: func() (*server, net.Conn) {
				mockLogger := internal_mock.NewMockLogger(ctrl)
				mockListener := net_mock.NewMockListener(ctrl)
				mockConn := net_mock.NewMockConn(ctrl)

				s := &server{
					logger:      mockLogger,
					listener:    mockListener,
					bufferSize:  1024,
					semaphore:   concurrency.NewSemaphore(1),
					idleTimeout: time.Second,
				}

				// Setup mocks for listener
				mockListener.EXPECT().Accept().Return(mockConn, nil).Times(1)
				mockListener.EXPECT().Accept().Return(nil, net.ErrClosed).Times(1)
				mockListener.EXPECT().Close().Return(nil).Times(1)

				// Setup mocks for connection
				mockConn.EXPECT().SetReadDeadline(gomock.Any()).Return(errors.New("some error")).Times(1)
				mockConn.EXPECT().Close().Return(nil).Times(1)

				// Setup mocks for logger
				mockLogger.EXPECT().Warn("failed to set read deadline", gomock.Any()).Times(1)

				return s, mockConn
			},
		},
		{
			name: "conn.Read count equal buffer size error",
			setupMocks: func() (*server, net.Conn) {
				mockLogger := internal_mock.NewMockLogger(ctrl)
				mockListener := net_mock.NewMockListener(ctrl)
				mockConn := net_mock.NewMockConn(ctrl)

				bufferSize := 1024
				s := &server{
					logger:     mockLogger,
					listener:   mockListener,
					bufferSize: bufferSize,
					semaphore:  concurrency.NewSemaphore(1),
				}

				// Setup mocks for listener
				mockListener.EXPECT().Accept().Return(mockConn, nil).Times(1)
				mockListener.EXPECT().Accept().Return(nil, net.ErrClosed).Times(1)
				mockListener.EXPECT().Close().Return(nil).Times(1)

				// Setup mocks for connection
				mockConn.EXPECT().Read(gomock.Any()).DoAndReturn(func(b []byte) (int, error) {
					return bufferSize, nil
				}).Times(1)
				mockConn.EXPECT().Close().Return(nil).Times(1)

				// Setup mocks for logger
				mockLogger.EXPECT().Warn("small buffer size", zap.Int("buffer_size", 1024)).Times(1)

				return s, mockConn
			},
		},
		{
			name: "conn.SetWriteDeadline returns an error",
			setupMocks: func() (*server, net.Conn) {
				mockLogger := internal_mock.NewMockLogger(ctrl)
				mockListener := net_mock.NewMockListener(ctrl)
				mockConn := net_mock.NewMockConn(ctrl)

				s := &server{
					logger:      mockLogger,
					listener:    mockListener,
					bufferSize:  1024,
					semaphore:   concurrency.NewSemaphore(1),
					idleTimeout: time.Second,
				}

				// Setup mocks for listener
				mockListener.EXPECT().Accept().Return(mockConn, nil).Times(1)
				mockListener.EXPECT().Accept().Return(nil, net.ErrClosed).Times(1)
				mockListener.EXPECT().Close().Return(nil).Times(1)

				// Setup mocks for connection
				mockConn.EXPECT().SetReadDeadline(gomock.Any()).Return(nil).Times(1)
				mockConn.EXPECT().SetWriteDeadline(gomock.Any()).Return(errors.New("write deadline error")).Times(1)
				mockConn.EXPECT().Read(gomock.Any()).Times(1)
				mockConn.EXPECT().Close().Times(1)

				// Setup mocks for logger
				mockLogger.EXPECT().Warn("failed to set read deadline", zap.Error(errors.New("write deadline error"))).Times(1)

				return s, mockConn
			},
		},
		{
			name: "conn.Write returns an error",
			setupMocks: func() (*server, net.Conn) {
				mockLogger := internal_mock.NewMockLogger(ctrl)
				mockListener := net_mock.NewMockListener(ctrl)
				mockConn := net_mock.NewMockConn(ctrl)

				s := &server{
					logger:     mockLogger,
					listener:   mockListener,
					bufferSize: 1024,
					semaphore:  concurrency.NewSemaphore(1),
				}

				// Setup mocks for listener
				mockListener.EXPECT().Accept().Return(mockConn, nil).Times(1)
				mockListener.EXPECT().Accept().Return(nil, net.ErrClosed).Times(1)
				mockListener.EXPECT().Close().Return(nil).Times(1)

				// Setup mocks for net.Addr
				mockAddr := net_mock.NewMockAddr(ctrl)
				mockAddr.EXPECT().String().Return("127.0.0.1:12345").Times(1)

				// Setup mocks for connection
				mockConn.EXPECT().Read(gomock.Any()).Times(1)
				mockConn.EXPECT().Write(gomock.Any()).Return(0, errors.New("write error")).Times(1)
				mockConn.EXPECT().RemoteAddr().Return(mockAddr).Times(1)
				mockConn.EXPECT().Close().Return(nil).Times(1)

				// Setup mocks for logger
				mockLogger.EXPECT().Warn(
					"failed to write data",
					zap.String("address", "127.0.0.1:12345"),
					zap.Error(errors.New("write error")),
				).Times(1)

				return s, mockConn
			},
		},
		{
			name: "conn.Close returns an error",
			setupMocks: func() (*server, net.Conn) {
				mockLogger := internal_mock.NewMockLogger(ctrl)
				mockListener := net_mock.NewMockListener(ctrl)
				mockConn := net_mock.NewMockConn(ctrl)

				s := &server{
					logger:      mockLogger,
					listener:    mockListener,
					bufferSize:  1024,
					semaphore:   concurrency.NewSemaphore(1),
					idleTimeout: time.Millisecond,
				}

				// Setup mocks for listener
				mockListener.EXPECT().Accept().Return(mockConn, nil).Times(1)
				mockListener.EXPECT().Accept().Return(nil, net.ErrClosed).Times(1)
				mockListener.EXPECT().Close().Return(nil).Times(1)

				// Setup mocks for connection
				mockConn.EXPECT().SetReadDeadline(gomock.Any()).Return(errors.New("some error")).Times(1) // break of cycle and move to defer
				mockConn.EXPECT().Close().Return(errors.New("connection close error")).Times(1)

				// Setup mocks for logger
				mockLogger.EXPECT().Warn("failed to set read deadline", gomock.Any()).Times(1)
				mockLogger.EXPECT().Warn("failed to close connection", zap.Error(errors.New("connection close error"))).Times(1)

				return s, mockConn
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			s, mockConn := tt.setupMocks()

			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()

			handler := func(ctx context.Context, data []byte) []byte {
				if tt.name == "Successful query handling" {
					assert.Equal(t, []byte("request"), data)
				}
				return []byte("response")
			}

			s.HandleQueries(ctx, handler)

			assert.NotNil(t, mockConn)
		})
	}
}

func TestHandleConnection_Timeout(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Setup mocks
	mockLogger := internal_mock.NewMockLogger(ctrl)
	mockConn := net_mock.NewMockConn(ctrl)

	server := &server{
		logger:      mockLogger,
		bufferSize:  1024,
		idleTimeout: time.Millisecond * 10,
	}

	mockLogger.EXPECT().Error(gomock.Any(), gomock.Any()).Times(1)

	mockConn.EXPECT().RemoteAddr().Return(nil).Times(1)
	mockConn.EXPECT().SetReadDeadline(gomock.Any()).Return(nil).Times(1)
	mockConn.EXPECT().Read(gomock.Any()).Return(0, errors.New("timeout")).Times(1)
	mockConn.EXPECT().Close().Return(nil).AnyTimes()

	handler := func(ctx context.Context, data []byte) []byte {
		return []byte("response")
	}

	server.handleConnection(context.Background(), mockConn, handler)
}
