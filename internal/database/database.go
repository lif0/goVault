package database

import (
	"context"
	"errors"
	"fmt"

	"goVault/internal"
	"goVault/internal/core/query"
	"goVault/internal/core/vault"
)

func NewDatabase(parserLayer parserLayer, vaultLayer vaultLayer, logger internal.Logger) (*Database, error) {
	if parserLayer == nil {
		return nil, errors.New("parser is invalid")
	}

	if vaultLayer == nil {
		return nil, errors.New("vault is invalid")
	}

	if logger == nil {
		return nil, errors.New("logger is invalid")
	}

	return &Database{
		parserLayer: parserLayer,
		vaultLayer:  vaultLayer,
		logger:      logger,
	}, nil
}

func (d *Database) HandleQuery(ctx context.Context, queryStr string) string {
	d.logger.Debug(fmt.Sprintf("handling query: query %s", queryStr))
	q, err := d.parserLayer.Transition(queryStr)
	if err != nil {
		return fmt.Sprintf("[error] %s", err.Error())
	}

	switch q.CommandID {
	case query.SET:
		return d.handleSetQuery(ctx, q)
	case query.GET:
		return d.handleGetQuery(ctx, q)
	case query.DEL:
		return d.handleDelQuery(ctx, q)
	}

	msg := fmt.Sprintf("compute layer is incorrect: command_id %v", q.CommandID)
	d.logger.Error(msg)

	return msg
}

func (d *Database) handleSetQuery(ctx context.Context, query *query.Query) string {
	arguments := query.Arguments
	if err := d.vaultLayer.Set(ctx, arguments[0], arguments[1]); err != nil {
		return fmt.Sprintf("[error] %s", err.Error())
	}

	return MsgDBOk
}

func (d *Database) handleGetQuery(ctx context.Context, query *query.Query) string {
	arguments := query.Arguments
	value, err := d.vaultLayer.Get(ctx, arguments[0])
	if err == vault.ErrVaultNotFound {
		return "[not found]"
	} else if err != nil {
		return fmt.Sprintf("%s %s", MsgDBError, err.Error())
	}

	return fmt.Sprintf("%s %s", MsgDBOk, value)
}

func (d *Database) handleDelQuery(ctx context.Context, query *query.Query) string {
	arguments := query.Arguments
	if err := d.vaultLayer.Del(ctx, arguments[0]); err != nil {
		return fmt.Sprintf("[error] %s", err.Error())
	}

	return MsgDBOk
}
