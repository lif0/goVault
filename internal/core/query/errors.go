package query

import "errors"

var (
	ErrQueryUnknowCommand = errors.New("query: unknow db command")
	ErrQueryUnknowState   = errors.New("query: fsm unknow state")

	ErrQueryInvalidArgumentSyntaxSet = errors.New("query: invalid argument syntax for the SET command")
	ErrQueryInvalidArgumentSyntaxGet = errors.New("query: invalid argument syntax for the SET command")
	ErrQueryInvalidArgumentSyntaxDel = errors.New("query: invalid argument syntax for the SET command")
)
