package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"syscall"

	"go.uber.org/zap"

	"goVault/internal/core/query"
	"goVault/internal/core/vault"
	"goVault/internal/core/vault/engine/in_memory"
	"goVault/internal/database"
)

func main() {
	logger, _ := zap.NewProduction()
	parser := query.NewParser()

	memEngine, err := in_memory.NewEngine()
	if err != nil {
		logger.Fatal("failed init engine")
	}

	vault, err := vault.NewVault(memEngine, logger)
	if err != nil {
		logger.Fatal("failed init vault")
	}

	db, err := database.NewDatabase(parser, vault, logger)
	if err != nil {
		logger.Fatal("failed init database")
	}

	//reader := bufio.NewReader(os.Stdin)
	reader := bufio.NewScanner(os.Stdin)

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

		response := db.HandleQuery(context.Background(), reader.Text())
		if errors.Is(err, syscall.EPIPE) {
			logger.Fatal("connection was closed", zap.Error(err))
		} else if err != nil {
			logger.Error("failed to send query", zap.Error(err))
		}

		fmt.Println(response)
	}
}
