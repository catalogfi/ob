// Code generated by MockGen. DO NOT EDIT.
// Source: ./swapper/bitcoin/client.go

// Package mock_bitcoin is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	btcec "github.com/btcsuite/btcd/btcec/v2"
	btcutil "github.com/btcsuite/btcd/btcutil"
	chaincfg "github.com/btcsuite/btcd/chaincfg"
	wire "github.com/btcsuite/btcd/wire"
	bitcoin "github.com/catalogfi/wbtc-garden/swapper/bitcoin"
	gomock "go.uber.org/mock/gomock"
)

// MockBitcoinClient is a mock of Client interface.
type MockBitcoinClient struct {
	ctrl     *gomock.Controller
	recorder *MockClientMockRecorder
}

// MockClientMockRecorder is the mock recorder for MockBitcoinClient.
type MockClientMockRecorder struct {
	mock *MockBitcoinClient
}

// NewMockBitcoinClient creates a new mock instance.
func NewMockBitcoinClient(ctrl *gomock.Controller) *MockBitcoinClient {
	mock := &MockBitcoinClient{ctrl: ctrl}
	mock.recorder = &MockClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBitcoinClient) EXPECT() *MockClientMockRecorder {
	return m.recorder
}

// CalculateRedeemFee mocks base method.
func (m *MockBitcoinClient) CalculateRedeemFee() (uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CalculateRedeemFee")
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CalculateRedeemFee indicates an expected call of CalculateRedeemFee.
func (mr *MockClientMockRecorder) CalculateRedeemFee() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CalculateRedeemFee", reflect.TypeOf((*MockBitcoinClient)(nil).CalculateRedeemFee))
}

// CalculateTransferFee mocks base method.
func (m *MockBitcoinClient) CalculateTransferFee(nInputs, nOutputs int, txVersion int32) (uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CalculateTransferFee", nInputs, nOutputs, txVersion)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CalculateTransferFee indicates an expected call of CalculateTransferFee.
func (mr *MockClientMockRecorder) CalculateTransferFee(nInputs, nOutputs, txVersion interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CalculateTransferFee", reflect.TypeOf((*MockBitcoinClient)(nil).CalculateTransferFee), nInputs, nOutputs, txVersion)
}

// GetConfirmations mocks base method.
func (m *MockBitcoinClient) GetConfirmations(txHash string) (uint64, uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetConfirmations", txHash)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(uint64)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetConfirmations indicates an expected call of GetConfirmations.
func (mr *MockClientMockRecorder) GetConfirmations(txHash interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetConfirmations", reflect.TypeOf((*MockBitcoinClient)(nil).GetConfirmations), txHash)
}

// GetFeeRates mocks base method.
func (m *MockBitcoinClient) GetFeeRates() (bitcoin.FeeRates, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFeeRates")
	ret0, _ := ret[0].(bitcoin.FeeRates)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetFeeRates indicates an expected call of GetFeeRates.
func (mr *MockClientMockRecorder) GetFeeRates() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFeeRates", reflect.TypeOf((*MockBitcoinClient)(nil).GetFeeRates))
}

// GetSpendingWitness mocks base method.
func (m *MockBitcoinClient) GetSpendingWitness(address btcutil.Address) ([]string, bitcoin.Transaction, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSpendingWitness", address)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(bitcoin.Transaction)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetSpendingWitness indicates an expected call of GetSpendingWitness.
func (mr *MockClientMockRecorder) GetSpendingWitness(address interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSpendingWitness", reflect.TypeOf((*MockBitcoinClient)(nil).GetSpendingWitness), address)
}

// GetTipBlockHeight mocks base method.
func (m *MockBitcoinClient) GetTipBlockHeight() (uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTipBlockHeight")
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTipBlockHeight indicates an expected call of GetTipBlockHeight.
func (mr *MockClientMockRecorder) GetTipBlockHeight() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTipBlockHeight", reflect.TypeOf((*MockBitcoinClient)(nil).GetTipBlockHeight))
}

// GetTx mocks base method.
func (m *MockBitcoinClient) GetTx(txid string) (bitcoin.Transaction, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTx", txid)
	ret0, _ := ret[0].(bitcoin.Transaction)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTx indicates an expected call of GetTx.
func (mr *MockClientMockRecorder) GetTx(txid interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTx", reflect.TypeOf((*MockBitcoinClient)(nil).GetTx), txid)
}

// GetUTXOs mocks base method.
func (m *MockBitcoinClient) GetUTXOs(address btcutil.Address, amount uint64) (bitcoin.UTXOs, uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUTXOs", address, amount)
	ret0, _ := ret[0].(bitcoin.UTXOs)
	ret1, _ := ret[1].(uint64)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetUTXOs indicates an expected call of GetUTXOs.
func (mr *MockClientMockRecorder) GetUTXOs(address, amount interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUTXOs", reflect.TypeOf((*MockBitcoinClient)(nil).GetUTXOs), address, amount)
}

// Net mocks base method.
func (m *MockBitcoinClient) Net() *chaincfg.Params {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Net")
	ret0, _ := ret[0].(*chaincfg.Params)
	return ret0
}

// Net indicates an expected call of Net.
func (mr *MockClientMockRecorder) Net() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Net", reflect.TypeOf((*MockBitcoinClient)(nil).Net))
}

// Send mocks base method.
func (m *MockBitcoinClient) Send(to btcutil.Address, amount uint64, from *btcec.PrivateKey) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Send", to, amount, from)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Send indicates an expected call of Send.
func (mr *MockClientMockRecorder) Send(to, amount, from interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Send", reflect.TypeOf((*MockBitcoinClient)(nil).Send), to, amount, from)
}

// Spend mocks base method.
func (m *MockBitcoinClient) Spend(script []byte, scriptSig wire.TxWitness, spender *btcec.PrivateKey, waitBlocks uint) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Spend", script, scriptSig, spender, waitBlocks)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Spend indicates an expected call of Spend.
func (mr *MockClientMockRecorder) Spend(script, scriptSig, spender, waitBlocks interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Spend", reflect.TypeOf((*MockBitcoinClient)(nil).Spend), script, scriptSig, spender, waitBlocks)
}
