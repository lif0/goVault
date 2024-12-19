package main

import (
	"bytes"
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"goVault/internal/configuration"
	"goVault/internal/initialization"
)

var (
	ConfigFileName = os.Getenv("CONFIG_FILE_NAME")
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg := getServerConfiguration()
	initializer, err := initialization.NewInitializer(cfg)
	if err != nil {
		log.Fatal(err)
	}

	if err = initializer.StartDatabase(ctx); err != nil {
		log.Fatal(err)
	}
}

func getServerConfiguration() *configuration.Config {
	cfg := &configuration.Config{}
	if ConfigFileName != "" {
		data, err := os.ReadFile(ConfigFileName)
		if err != nil {
			log.Fatal(err)
		}

		reader := bytes.NewReader(data)
		cfg, err = configuration.Load(reader)
		if err != nil {
			log.Fatal(err)
		}
	}

	return cfg
}
