// Code generated by MockGen. DO NOT EDIT.
// Source: ../swapper/swapper.go

// Package mock_swapper is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	model "github.com/catalogfi/wbtc-garden/model"
	gomock "go.uber.org/mock/gomock"
)

// MockInitiatorSwap is a mock of InitiatorSwap interface.
type MockInitiatorSwap struct {
	ctrl     *gomock.Controller
	recorder *MockInitiatorSwapMockRecorder
}

// MockInitiatorSwapMockRecorder is the mock recorder for MockInitiatorSwap.
type MockInitiatorSwapMockRecorder struct {
	mock *MockInitiatorSwap
}

// NewMockInitiatorSwap creates a new mock instance.
func NewMockInitiatorSwap(ctrl *gomock.Controller) *MockInitiatorSwap {
	mock := &MockInitiatorSwap{ctrl: ctrl}
	mock.recorder = &MockInitiatorSwapMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockInitiatorSwap) EXPECT() *MockInitiatorSwapMockRecorder {
	return m.recorder
}

// Expired mocks base method.
func (m *MockInitiatorSwap) Expired() (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Expired")
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Expired indicates an expected call of Expired.
func (mr *MockInitiatorSwapMockRecorder) Expired() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Expired", reflect.TypeOf((*MockInitiatorSwap)(nil).Expired))
}

// Initiate mocks base method.
func (m *MockInitiatorSwap) Initiate() (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Initiate")
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Initiate indicates an expected call of Initiate.
func (mr *MockInitiatorSwapMockRecorder) Initiate() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Initiate", reflect.TypeOf((*MockInitiatorSwap)(nil).Initiate))
}

// IsRedeemed mocks base method.
func (m *MockInitiatorSwap) IsRedeemed() (bool, []byte, string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsRedeemed")
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].([]byte)
	ret2, _ := ret[2].(string)
	ret3, _ := ret[3].(error)
	return ret0, ret1, ret2, ret3
}

// IsRedeemed indicates an expected call of IsRedeemed.
func (mr *MockInitiatorSwapMockRecorder) IsRedeemed() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsRedeemed", reflect.TypeOf((*MockInitiatorSwap)(nil).IsRedeemed))
}

// Refund mocks base method.
func (m *MockInitiatorSwap) Refund() (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Refund")
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Refund indicates an expected call of Refund.
func (mr *MockInitiatorSwapMockRecorder) Refund() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Refund", reflect.TypeOf((*MockInitiatorSwap)(nil).Refund))
}

// WaitForRedeem mocks base method.
func (m *MockInitiatorSwap) WaitForRedeem() ([]byte, string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WaitForRedeem")
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(string)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// WaitForRedeem indicates an expected call of WaitForRedeem.
func (mr *MockInitiatorSwapMockRecorder) WaitForRedeem() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WaitForRedeem", reflect.TypeOf((*MockInitiatorSwap)(nil).WaitForRedeem))
}

// MockRedeemerSwap is a mock of RedeemerSwap interface.
type MockRedeemerSwap struct {
	ctrl     *gomock.Controller
	recorder *MockRedeemerSwapMockRecorder
}

// MockRedeemerSwapMockRecorder is the mock recorder for MockRedeemerSwap.
type MockRedeemerSwapMockRecorder struct {
	mock *MockRedeemerSwap
}

// NewMockRedeemerSwap creates a new mock instance.
func NewMockRedeemerSwap(ctrl *gomock.Controller) *MockRedeemerSwap {
	mock := &MockRedeemerSwap{ctrl: ctrl}
	mock.recorder = &MockRedeemerSwapMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRedeemerSwap) EXPECT() *MockRedeemerSwapMockRecorder {
	return m.recorder
}

// IsInitiated mocks base method.
func (m *MockRedeemerSwap) IsInitiated() (bool, string, uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsInitiated")
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(string)
	ret2, _ := ret[2].(uint64)
	ret3, _ := ret[3].(error)
	return ret0, ret1, ret2, ret3
}

// IsInitiated indicates an expected call of IsInitiated.
func (mr *MockRedeemerSwapMockRecorder) IsInitiated() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsInitiated", reflect.TypeOf((*MockRedeemerSwap)(nil).IsInitiated))
}

// Redeem mocks base method.
func (m *MockRedeemerSwap) Redeem(secret []byte) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Redeem", secret)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Redeem indicates an expected call of Redeem.
func (mr *MockRedeemerSwapMockRecorder) Redeem(secret interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Redeem", reflect.TypeOf((*MockRedeemerSwap)(nil).Redeem), secret)
}

// WaitForInitiate mocks base method.
func (m *MockRedeemerSwap) WaitForInitiate() (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WaitForInitiate")
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// WaitForInitiate indicates an expected call of WaitForInitiate.
func (mr *MockRedeemerSwapMockRecorder) WaitForInitiate() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WaitForInitiate", reflect.TypeOf((*MockRedeemerSwap)(nil).WaitForInitiate))
}

// MockWatcher is a mock of Watcher interface.
type MockWatcher struct {
	ctrl     *gomock.Controller
	recorder *MockWatcherMockRecorder
}

// MockWatcherMockRecorder is the mock recorder for MockWatcher.
type MockWatcherMockRecorder struct {
	mock *MockWatcher
}

// NewMockWatcher creates a new mock instance.
func NewMockWatcher(ctrl *gomock.Controller) *MockWatcher {
	mock := &MockWatcher{ctrl: ctrl}
	mock.recorder = &MockWatcherMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockWatcher) EXPECT() *MockWatcherMockRecorder {
	return m.recorder
}

// Expired mocks base method.
func (m *MockWatcher) Expired() (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Expired")
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Expired indicates an expected call of Expired.
func (mr *MockWatcherMockRecorder) Expired() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Expired", reflect.TypeOf((*MockWatcher)(nil).Expired))
}

// Identifier mocks base method.
func (m *MockWatcher) Identifier() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Identifier")
	ret0, _ := ret[0].(string)
	return ret0
}

// Identifier indicates an expected call of Identifier.
func (mr *MockWatcherMockRecorder) Identifier() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Identifier", reflect.TypeOf((*MockWatcher)(nil).Identifier))
}

// IsDetected mocks base method.
func (m *MockWatcher) IsDetected() (bool, string, string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsDetected")
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(string)
	ret2, _ := ret[2].(string)
	ret3, _ := ret[3].(error)
	return ret0, ret1, ret2, ret3
}

// IsDetected indicates an expected call of IsDetected.
func (mr *MockWatcherMockRecorder) IsDetected() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsDetected", reflect.TypeOf((*MockWatcher)(nil).IsDetected))
}

// IsInitiated mocks base method.
func (m *MockWatcher) IsInitiated() (bool, string, map[string]model.Chain, uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsInitiated")
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(string)
	ret2, _ := ret[2].(map[string]model.Chain)
	ret3, _ := ret[3].(uint64)
	ret4, _ := ret[4].(error)
	return ret0, ret1, ret2, ret3, ret4
}

// IsInitiated indicates an expected call of IsInitiated.
func (mr *MockWatcherMockRecorder) IsInitiated() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsInitiated", reflect.TypeOf((*MockWatcher)(nil).IsInitiated))
}

// IsRedeemed mocks base method.
func (m *MockWatcher) IsRedeemed() (bool, []byte, string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsRedeemed")
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].([]byte)
	ret2, _ := ret[2].(string)
	ret3, _ := ret[3].(error)
	return ret0, ret1, ret2, ret3
}

// IsRedeemed indicates an expected call of IsRedeemed.
func (mr *MockWatcherMockRecorder) IsRedeemed() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsRedeemed", reflect.TypeOf((*MockWatcher)(nil).IsRedeemed))
}

// IsRefunded mocks base method.
func (m *MockWatcher) IsRefunded() (bool, string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsRefunded")
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(string)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// IsRefunded indicates an expected call of IsRefunded.
func (mr *MockWatcherMockRecorder) IsRefunded() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsRefunded", reflect.TypeOf((*MockWatcher)(nil).IsRefunded))
}

// Status mocks base method.
func (m *MockWatcher) Status(initiateTxHash string) (uint64, uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Status", initiateTxHash)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(uint64)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Status indicates an expected call of Status.
func (mr *MockWatcherMockRecorder) Status(initiateTxHash interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Status", reflect.TypeOf((*MockWatcher)(nil).Status), initiateTxHash)
}
