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
func (m *MockVault) Del(arg0 context.Context, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Del", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Del indicates an expected call of Del.
func (mr *MockVaultMockRecorder) Del(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Del", reflect.TypeOf((*MockVault)(nil).Del), arg0, arg1)
}

// Get mocks base method.
func (m *MockVault) Get(arg0 context.Context, arg1 string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", arg0, arg1)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockVaultMockRecorder) Get(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockVault)(nil).Get), arg0, arg1)
}

// Set mocks base method.
func (m *MockVault) Set(arg0 context.Context, arg1, arg2 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Set", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// Set indicates an expected call of Set.
func (mr *MockVaultMockRecorder) Set(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Set", reflect.TypeOf((*MockVault)(nil).Set), arg0, arg1, arg2)
}