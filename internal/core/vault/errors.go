package vault

import "errors"

var (
	ErrVaultNotFound       = errors.New("vault: key not found")
	ErrVaultUnknownCommand = errors.New("vault: unknown db command")
)
