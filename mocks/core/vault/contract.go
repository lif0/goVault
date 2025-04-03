// Code generated by MockGen. DO NOT EDIT.
// Source: goVault/internal/core/vault (interfaces: Vault)
//
// Generated by this command:
//
//	mockgen -destination ./../../../mocks/core/vault/contract.go -package vault_mock . Vault
//

// Package vault_mock is a generated GoMock package.
package vault_mock

import (
	context "context"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockVault is a mock of Vault interface.
type MockVault struct {
	ctrl     *gomock.Controller
	recorder *MockVaultMockRecorder
	isgomock struct{}
}

// MockVaultMockRecorder is the mock recorder for MockVault.
type MockVaultMockRecorder struct {
	mock *MockVault
}

// NewMockVault creates a new mock instance.
func NewMockVault(ctrl *gomock.Controller) *MockVault {
	mock := &MockVault{ctrl: ctrl}
	mock.recorder = &MockVaultMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockVault) EXPECT() *MockVaultMockRecorder {
	return m.recorder
}

// Del mocks base method.
func (m *MockVault) Del(ctx context.Context, key string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Del", ctx, key)
	ret0, _ := ret[0].(error)
	return ret0
}

// Del indicates an expected call of Del.
func (mr *MockVaultMockRecorder) Del(ctx, key any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Del", reflect.TypeOf((*MockVault)(nil).Del), ctx, key)
}

// Get mocks base method.
func (m *MockVault) Get(ctx context.Context, key string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, key)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockVaultMockRecorder) Get(ctx, key any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockVault)(nil).Get), ctx, key)
}

// Set mocks base method.
func (m *MockVault) Set(ctx context.Context, key, value string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Set", ctx, key, value)
	ret0, _ := ret[0].(error)
	return ret0
}

// Set indicates an expected call of Set.
func (mr *MockVaultMockRecorder) Set(ctx, key, value any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Set", reflect.TypeOf((*MockVault)(nil).Set), ctx, key, value)
}
