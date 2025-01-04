package in_memory

import (
	"context"
	"testing"

	"go.uber.org/mock/gomock"

	internal_mock "goVault/mocks"
)

func TestNewEngine(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := internal_mock.NewMockLogger(ctrl)
	eng, err := NewEngine(logger)

	if err != nil {
		t.Fatalf("expected no error while creating engine, got %v", err)
	}
	if eng == nil {
		t.Fatalf("expected engine instance to be created, got nil")
	}
}

func TestService_Set(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := internal_mock.NewMockLogger(ctrl)
	eng, _ := NewEngine(logger)

	ctx := context.Background()

	// Expect logger to log a debug message for "Set" operation
	logger.EXPECT().Debug("successfull set query").Times(1)

	eng.Set(ctx, "key1", "value1")

	// Validate that the key-value pair was set correctly
	value, ok := eng.(*service).vault.Get("key1")
	if !ok {
		t.Fatalf("expected key to be present")
	}
	if value != "value1" {
		t.Fatalf("expected value to be 'value1', got %v", value)
	}
}

func TestService_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := internal_mock.NewMockLogger(ctrl)
	logger.EXPECT().Debug(gomock.Any()).AnyTimes()
	eng, _ := NewEngine(logger)

	ctx := context.Background()
	eng.Set(ctx, "key1", "value1")

	value, ok := eng.Get(ctx, "key1")
	if !ok {
		t.Fatalf("expected key to be present")
	}
	if value != "value1" {
		t.Fatalf("expected value to be 'value1', got %v", value)
	}
}

func TestService_Del(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := internal_mock.NewMockLogger(ctrl)
	logger.EXPECT().Debug(gomock.Any()).AnyTimes()
	eng, _ := NewEngine(logger)

	ctx := context.Background()
	eng.Set(ctx, "key1", "value1")

	eng.Del(ctx, "key1")

	// Validate that the key was deleted
	_, ok := eng.Get(ctx, "key1")
	if ok {
		t.Fatalf("expected key to be absent after deletion")
	}
}
