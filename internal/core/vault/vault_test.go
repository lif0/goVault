package vault

import (
	"context"
	"errors"
	"testing"

	"go.uber.org/mock/gomock"

	internal_mock "goVault/mocks"
	engine_mock "goVault/mocks/core/vault/engine/in_memory"
)

func TestNewVault(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := internal_mock.NewMockLogger(ctrl)
	engine := engine_mock.NewMockEngine(ctrl)

	// Valid inputs
	v, err := NewVault(engine, logger)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if v == nil {
		t.Fatalf("expected vault instance, got nil")
	}

	// Invalid engine
	_, err = NewVault(nil, logger)
	if err == nil || err.Error() != "engine is invalid" {
		t.Fatalf("expected 'engine is invalid' error, got %v", err)
	}

	// Invalid logger
	_, err = NewVault(engine, nil)
	if err == nil || err.Error() != "logger is invalid" {
		t.Fatalf("expected 'logger is invalid' error, got %v", err)
	}
}

func TestVault_Set(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := internal_mock.NewMockLogger(ctrl)
	engine := engine_mock.NewMockEngine(ctrl)
	v, _ := NewVault(engine, logger)

	ctx := context.Background()

	// Expect engine.Set to be called with correct arguments
	engine.EXPECT().Set(ctx, "key1", "value1").Times(1)

	err := v.Set(ctx, "key1", "value1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Test with cancelled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err = v.Set(ctx, "key1", "value1")
	if err == nil || !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled error, got %v", err)
	}
}

func TestVault_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := internal_mock.NewMockLogger(ctrl)
	engine := engine_mock.NewMockEngine(ctrl)
	v, _ := NewVault(engine, logger)

	ctx := context.Background()

	// Expect engine.Get to return a valid value
	engine.EXPECT().Get(ctx, "key1").Return("value1", true).Times(1)

	value, err := v.Get(ctx, "key1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if value != "value1" {
		t.Fatalf("expected value 'value1', got %v", value)
	}

	// Expect engine.Get to return not found
	engine.EXPECT().Get(ctx, "key2").Return("", false).Times(1)

	_, err = v.Get(ctx, "key2")
	if err == nil || err != ErrVaultNotFound {
		t.Fatalf("expected ErrVaultNotFound, got %v", err)
	}

	// Test with cancelled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	value, err = v.Get(ctx, "key1")
	if err == nil || !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled error, got %v", err)
	}
	if value != "" {
		t.Fatalf("expected empty value, got %v", value)
	}
}

func TestVault_Del(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := internal_mock.NewMockLogger(ctrl)
	engine := engine_mock.NewMockEngine(ctrl)
	v, _ := NewVault(engine, logger)

	ctx := context.Background()

	// Expect engine.Del to be called
	engine.EXPECT().Del(ctx, "key1").Times(1)

	err := v.Del(ctx, "key1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Test with cancelled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err = v.Del(ctx, "key1")
	if err == nil || !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled error, got %v", err)
	}
}
