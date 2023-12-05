// Code generated by MockGen. DO NOT EDIT.
// Source: ./swapper/bitcoin/store.go

// Package mock_bitcoin is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	bitcoin "github.com/catalogfi/orderbook/swapper/bitcoin"
	gomock "go.uber.org/mock/gomock"
)

// MockBTCStore is a mock of Store interface.
type MockBTCStore struct {
	ctrl     *gomock.Controller
	recorder *MockBTCStoreMockRecorder
}

// MockBTCStoreMockRecorder is the mock recorder for MockBTCStore.
type MockBTCStoreMockRecorder struct {
	mock *MockBTCStore
}

// NewMockBTCStore creates a new mock instance.
func NewMockBTCStore(ctrl *gomock.Controller) *MockBTCStore {
	mock := &MockBTCStore{ctrl: ctrl}
	mock.recorder = &MockBTCStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBTCStore) EXPECT() *MockBTCStoreMockRecorder {
	return m.recorder
}

// PutSecret mocks base method.
func (m *MockBTCStore) PutSecret(pubkey, secret string, status bitcoin.IwStatus, code uint32) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PutSecret", pubkey, secret, status, code)
	ret0, _ := ret[0].(error)
	return ret0
}

// PutSecret indicates an expected call of PutSecret.
func (mr *MockBTCStoreMockRecorder) PutSecret(pubkey, secret, status, code interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PutSecret", reflect.TypeOf((*MockBTCStore)(nil).PutSecret), pubkey, secret, status, code)
}

// PutStatus mocks base method.
func (m *MockBTCStore) PutStatus(pubkey string, code uint32, status bitcoin.IwStatus) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PutStatus", pubkey, code, status)
	ret0, _ := ret[0].(error)
	return ret0
}

// PutStatus indicates an expected call of PutStatus.
func (mr *MockBTCStoreMockRecorder) PutStatus(pubkey, code, status interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PutStatus", reflect.TypeOf((*MockBTCStore)(nil).PutStatus), pubkey, code, status)
}

// Secret mocks base method.
func (m *MockBTCStore) Secret(pubkey string, code uint32) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Secret", pubkey, code)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Secret indicates an expected call of Secret.
func (mr *MockBTCStoreMockRecorder) Secret(pubkey, code interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Secret", reflect.TypeOf((*MockBTCStore)(nil).Secret), pubkey, code)
}

// Status mocks base method.
func (m *MockBTCStore) Status(pubkey string, code uint32) (bitcoin.IwStatus, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Status", pubkey, code)
	ret0, _ := ret[0].(bitcoin.IwStatus)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Status indicates an expected call of Status.
func (mr *MockBTCStoreMockRecorder) Status(pubkey, code interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Status", reflect.TypeOf((*MockBTCStore)(nil).Status), pubkey, code)
}
