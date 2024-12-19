package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"syscall"
	"time"

	"go.uber.org/zap"

	"goVault/client"
	"goVault/internal"
	"goVault/internal/pkg/unit"
)

type ClientArguments struct {
	Address        string
	IdleTimeout    time.Duration
	MaxMessageSize uint
}

func main() {
	logger, _ := zap.NewProduction()
	reader := bufio.NewScanner(os.Stdin)

	// options
	args := grabArguments(logger)
	fmt.Println(args.Address)

	var options []client.TCPClientOption
	options = append(options, client.WithClientIdleTimeout(args.IdleTimeout))
	options = append(options, client.WithClientBufferSize(args.MaxMessageSize))

	client, err := client.NewTCPClient(args.Address, options...)
	if err != nil {
		logger.Fatal("failed to connect with server", zap.Error(err))
	}

	// cycle
	for {
		fmt.Print("[goVault] > ")
		if !reader.Scan() { // scan console
			break
		}

		if errors.Is(err, syscall.EPIPE) {
			logger.Fatal("connection was closed", zap.Error(err))
		} else if err != nil {
			logger.Error("failed to read query", zap.Error(err))
		}

		response, err := client.Send(reader.Bytes())
		if errors.Is(err, syscall.EPIPE) {
			logger.Fatal("connection was closed", zap.Error(err))
		} else if err != nil {
			logger.Error("failed to send query", zap.Error(err))
		}

		fmt.Println(string(response))
	}
}

func grabArguments(logger internal.Logger) *ClientArguments {
	address := flag.String("address", "localhost:3231", "Address of the spider")
	idleTimeout := flag.Duration("idle_timeout", time.Minute, "Idle timeout for connection")
	maxMessageSizeStr := flag.String("max_message_size", "4KB", "Max message size for connection")
	flag.Parse()

	app_args := ClientArguments{}

	if address != nil && len(*address) > 0 {
		app_args.Address = *address
	}

	if idleTimeout != nil {
		app_args.IdleTimeout = *idleTimeout
	}

	if maxMessageSizeStr != nil && len(*maxMessageSizeStr) > 0 {
		maxMessageSize, err := unit.ParseDigitalStorage(*maxMessageSizeStr)
		if err != nil {
			logger.Fatal("failed to parse max message size", zap.Error(err))
		}

		app_args.MaxMessageSize = uint(maxMessageSize)
	}

	return &app_args
}
