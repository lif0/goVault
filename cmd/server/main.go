package main

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"

	"goVault/internal/configuration"
	"goVault/internal/initialization"
)

var (
	ConfigPath = os.Getenv("CONFIG_PATH")
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg := getServerConfiguration()

	if cfg == nil {
		log.Default().Println("Config is empty")
	} else {
		json, _ := json.MarshalIndent(cfg, "", "  ")
		log.Default().Printf("Config:\n%s\n", string(json))
	}

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
	if ConfigPath != "" {
		data, err := os.ReadFile(ConfigPath)
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
