package database

import "goVault/internal"

type Database struct {
	parserLayer parserLayer
	vaultLayer  vaultLayer
	logger      internal.Logger
}
