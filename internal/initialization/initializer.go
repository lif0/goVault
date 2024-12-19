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
	"goVault/internal/database"
	"goVault/internal/network"
)

type Initializer struct {
	engine engine.Engine
	logger internal.Logger
	server *network.TCPServer
}

func NewInitializer(cfg *configuration.Config) (*Initializer, error) {
	if cfg == nil {
		return nil, errors.New("failed to initialize: config is invalid")
	}

	logger, err := CreateLogger(cfg.Logging)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize logger: %w", err)
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

	vault, err := vault.NewVault(i.engine, i.logger)
	if err != nil {
		return err
	}

	database, err := database.NewDatabase(parser, vault, i.logger)
	if err != nil {
		return err
	}

	group, groupCtx := errgroup.WithContext(ctx)
	group.Go(func() error {
		i.server.HandleQueries(groupCtx, func(ctx context.Context, query []byte) []byte {
			response := database.HandleQuery(ctx, string(query))
			return []byte(response)
		})

		return nil
	})

	return group.Wait()
}
