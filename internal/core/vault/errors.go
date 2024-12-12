package vault

import "errors"

var (
	ErrVaultNotFound      = errors.New("vault: key not found")
	ErrVaultUnknowCommand = errors.New("vault: unknow db command")
)
