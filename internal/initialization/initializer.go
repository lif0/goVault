package initialization

import (
	"context"
	"errors"
	"fmt"

	"golang.org/x/sync/errgroup"

	"goVault/internal"
	"goVault/internal/configuration"
	"goVault/internal/core/query"
	"goVault/internal/core/vault"
	"goVault/internal/core/vault/engine"
	"goVault/internal/core/vault/wal"
	"goVault/internal/database"
	"goVault/internal/network"
)

type Initializer struct {
	logger internal.Logger
	engine engine.Engine
	wal    wal.WAL
	server network.TCPServer
}

func NewInitializer(cfg *configuration.Config) (*Initializer, error) {
	if cfg == nil {
		return nil, errors.New("failed to initialize: config is invalid")
	}

	logger, err := CreateLogger(cfg.Logging)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize logger: %w", err)
	}

	wal, err := CreateWAL(cfg.WAL, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize wal: %w", err)
	}

	engine, err := CreateEngine(cfg.Engine, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize engine: %w", err)
	}

	server, err := CreateNetwork(cfg.Network, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize network: %w", err)
	}

	initializer := &Initializer{
		wal:    wal,
		engine: engine,
		server: server,
		logger: logger,
	}

	return initializer, nil
}

func (i *Initializer) StartDatabase(ctx context.Context) error {
	parser, err := query.NewParser(i.logger)
	if err != nil {
		return err
	}

	var options []vault.VaultOption
	if i.wal != nil {
		options = append(options, vault.WithWAL(i.wal))
	}

	vault, err := vault.NewVault(i.engine, i.wal, i.logger, options...)
	if err != nil {
		return err
	}

	database, err := database.NewDatabase(parser, vault, i.logger)
	if err != nil {
		return err
	}

	group, groupCtx := errgroup.WithContext(ctx)

	if i.wal != nil {
		group.Go(func() error {
			i.wal.Start(groupCtx)
			return nil
		})
	}

	group.Go(func() error {
		i.server.HandleQueries(groupCtx, func(ctx context.Context, query []byte) []byte {
			response := database.HandleQuery(ctx, string(query))
			return []byte(response)
		})

		return nil
	})

	return group.Wait()
}
