package vault

import (
	"context"
	"errors"

	"goVault/internal"
)

type vault struct {
	engine Engine
	logger internal.Logger
}

func NewVault(e Engine, l internal.Logger) (Vault, error) {
	if e == nil {
		return nil, errors.New("engine is invalid")
	}

	if l == nil {
		return nil, errors.New("logger is invalid")
	}

	v := vault{
		engine: e,
		logger: l,
	}

	return &v, nil
}

func (v *vault) Set(ctx context.Context, key, value string) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	v.engine.Set(ctx, key, value)
	return nil
}

func (v *vault) Get(ctx context.Context, key string) (string, error) {
	if ctx.Err() != nil {
		return "", ctx.Err()
	}

	value, found := v.engine.Get(ctx, key)
	if !found {
		return "", ErrVaultNotFound
	}

	return value, nil
}

func (v *vault) Del(ctx context.Context, key string) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	v.engine.Del(ctx, key)
	return nil
}

// func (vp *vaultProcessor) ProcessQuery(rawQuery string) Result {
// 	result := Result{}
// 	query, err := vp.parser.Transition(rawQuery)
// 	if err != nil {
// 		result.Error = err
// 		return result
// 	}

// 	return vp.processQuery(query)
// }

// func (vp *vaultProcessor) processQuery(q *query.Query) Result {
// 	result := Result{}

// 	q_ := pointer.ValueOf(q)

// 	switch q_.CommandID {
// 	case query.SET:
// 		key := q_.Arguments[0]
// 		value := q_.Arguments[1]

// 		vp.engine.Set(key, value)
// 		result.Success = true

// 	case query.GET:
// 		key := q_.Arguments[0]

// 		raw, err := vp.engine.Get(key)
// 		result.Success = true
// 		result.RawValue = raw
// 		result.Error = err

// 	case query.DEL:
// 		key := q_.Arguments[0]

// 		vp.engine.Del(key)
// 		result.Success = true

// 	default:
// 		result.Error = ErrVaultUnknowCommand
// 	}

// 	return result
// }
