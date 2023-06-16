// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package AtomicSwap

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// AtomicSwapMetaData contains all meta data concerning the AtomicSwap contract.
var AtomicSwapMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_redeemer\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_refunder\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"_secretHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"_expiry\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"_secret\",\"type\":\"string\"}],\"name\":\"Redeemed\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"_secret\",\"type\":\"bytes\"}],\"name\":\"redeem\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"}],\"name\":\"refund\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x6101006040523480156200001257600080fd5b5060405162000cbf38038062000cbf83398181016040528101906200003891906200019a565b8373ffffffffffffffffffffffffffffffffffffffff1660808173ffffffffffffffffffffffffffffffffffffffff16815250508273ffffffffffffffffffffffffffffffffffffffff1660a08173ffffffffffffffffffffffffffffffffffffffff16815250508160c081815250508060e08181525050505050506200020c565b600080fd5b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b6000620000ec82620000bf565b9050919050565b620000fe81620000df565b81146200010a57600080fd5b50565b6000815190506200011e81620000f3565b92915050565b6000819050919050565b620001398162000124565b81146200014557600080fd5b50565b60008151905062000159816200012e565b92915050565b6000819050919050565b62000174816200015f565b81146200018057600080fd5b50565b600081519050620001948162000169565b92915050565b60008060008060808587031215620001b757620001b6620000ba565b5b6000620001c7878288016200010d565b9450506020620001da878288016200010d565b9350506040620001ed8782880162000148565b9250506060620002008782880162000183565b91505092959194509250565b60805160a05160c05160e051610a7a620002456000396000610155015260006075015260006101850152600061012b0152610a7a6000f3fe608060405234801561001057600080fd5b50600436106100365760003560e01c80633841185e1461003b578063fa89401a14610057575b600080fd5b6100556004803603810190610050919061060b565b610073565b005b610071600480360381019061006c9190610667565b610153565b005b7f00000000000000000000000000000000000000000000000000000000000000006002826040516100a49190610705565b602060405180830381855afa1580156100c1573d6000803e3d6000fd5b5050506040513d601f19601f820116820180604052508101906100e49190610752565b146100ee57600080fd5b7f434fa9d5634e23e27249731c499be7ac2854b2717ff296496d611e3157f036428160405161011d91906107d4565b60405180910390a161014f827f00000000000000000000000000000000000000000000000000000000000000006101ac565b5050565b7f0000000000000000000000000000000000000000000000000000000000000000431161017f57600080fd5b6101a9817f00000000000000000000000000000000000000000000000000000000000000006101ac565b50565b6000808373ffffffffffffffffffffffffffffffffffffffff166370a08231306040516024016101dc9190610805565b6040516020818303038152906040529060e01b6020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff838183161783525050505060405161022a9190610705565b6000604051808303816000865af19150503d8060008114610267576040519150601f19603f3d011682016040523d82523d6000602084013e61026c565b606091505b5091509150816102b1576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016102a890610892565b60405180910390fd5b6000818060200190518101906102c791906108e8565b9050600081111561043a576000808673ffffffffffffffffffffffffffffffffffffffff1663a9059cbb8785604051602401610304929190610924565b6040516020818303038152906040529060e01b6020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff83818316178352505050506040516103529190610705565b6000604051808303816000865af19150503d806000811461038f576040519150601f19603f3d011682016040523d82523d6000602084013e610394565b606091505b5091509150816103d9576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016103d0906109bf565b60405180910390fd5b60008151111561043757808060200190518101906103f79190610a17565b610436576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161042d906109bf565b60405180910390fd5b5b50505b8373ffffffffffffffffffffffffffffffffffffffff16ff5b6000604051905090565b600080fd5b600080fd5b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b600061049282610467565b9050919050565b6104a281610487565b81146104ad57600080fd5b50565b6000813590506104bf81610499565b92915050565b600080fd5b600080fd5b6000601f19601f8301169050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b610518826104cf565b810181811067ffffffffffffffff82111715610537576105366104e0565b5b80604052505050565b600061054a610453565b9050610556828261050f565b919050565b600067ffffffffffffffff821115610576576105756104e0565b5b61057f826104cf565b9050602081019050919050565b82818337600083830152505050565b60006105ae6105a98461055b565b610540565b9050828152602081018484840111156105ca576105c96104ca565b5b6105d584828561058c565b509392505050565b600082601f8301126105f2576105f16104c5565b5b813561060284826020860161059b565b91505092915050565b600080604083850312156106225761062161045d565b5b6000610630858286016104b0565b925050602083013567ffffffffffffffff81111561065157610650610462565b5b61065d858286016105dd565b9150509250929050565b60006020828403121561067d5761067c61045d565b5b600061068b848285016104b0565b91505092915050565b600081519050919050565b600081905092915050565b60005b838110156106c85780820151818401526020810190506106ad565b60008484015250505050565b60006106df82610694565b6106e9818561069f565b93506106f98185602086016106aa565b80840191505092915050565b600061071182846106d4565b915081905092915050565b6000819050919050565b61072f8161071c565b811461073a57600080fd5b50565b60008151905061074c81610726565b92915050565b6000602082840312156107685761076761045d565b5b60006107768482850161073d565b91505092915050565b600081519050919050565b600082825260208201905092915050565b60006107a68261077f565b6107b0818561078a565b93506107c08185602086016106aa565b6107c9816104cf565b840191505092915050565b600060208201905081810360008301526107ee818461079b565b905092915050565b6107ff81610487565b82525050565b600060208201905061081a60008301846107f6565b92915050565b7f41746f6d6963537761703a2045524332302062616c616e63654f66206469642060008201527f6e6f742073756363656564000000000000000000000000000000000000000000602082015250565b600061087c602b8361078a565b915061088782610820565b604082019050919050565b600060208201905081810360008301526108ab8161086f565b9050919050565b6000819050919050565b6108c5816108b2565b81146108d057600080fd5b50565b6000815190506108e2816108bc565b92915050565b6000602082840312156108fe576108fd61045d565b5b600061090c848285016108d3565b91505092915050565b61091e816108b2565b82525050565b600060408201905061093960008301856107f6565b6109466020830184610915565b9392505050565b7f41746f6d6963537761703a204552433230207472616e7366657220646964206e60008201527f6f74207375636365656400000000000000000000000000000000000000000000602082015250565b60006109a9602a8361078a565b91506109b48261094d565b604082019050919050565b600060208201905081810360008301526109d88161099c565b9050919050565b60008115159050919050565b6109f4816109df565b81146109ff57600080fd5b50565b600081519050610a11816109eb565b92915050565b600060208284031215610a2d57610a2c61045d565b5b6000610a3b84828501610a02565b9150509291505056fea264697066735822122021904c7f19b9b2114a5be2cc8b3664e601c6528e5e8c516a18d3b919bb32520b64736f6c63430008130033",
}

// AtomicSwapABI is the input ABI used to generate the binding from.
// Deprecated: Use AtomicSwapMetaData.ABI instead.
var AtomicSwapABI = AtomicSwapMetaData.ABI

// AtomicSwapBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use AtomicSwapMetaData.Bin instead.
var AtomicSwapBin = AtomicSwapMetaData.Bin

// DeployAtomicSwap deploys a new Ethereum contract, binding an instance of AtomicSwap to it.
func DeployAtomicSwap(auth *bind.TransactOpts, backend bind.ContractBackend, _redeemer common.Address, _refunder common.Address, _secretHash [32]byte, _expiry *big.Int) (common.Address, *types.Transaction, *AtomicSwap, error) {
	parsed, err := AtomicSwapMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(AtomicSwapBin), backend, _redeemer, _refunder, _secretHash, _expiry)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &AtomicSwap{AtomicSwapCaller: AtomicSwapCaller{contract: contract}, AtomicSwapTransactor: AtomicSwapTransactor{contract: contract}, AtomicSwapFilterer: AtomicSwapFilterer{contract: contract}}, nil
}

// AtomicSwap is an auto generated Go binding around an Ethereum contract.
type AtomicSwap struct {
	AtomicSwapCaller     // Read-only binding to the contract
	AtomicSwapTransactor // Write-only binding to the contract
	AtomicSwapFilterer   // Log filterer for contract events
}

// AtomicSwapCaller is an auto generated read-only Go binding around an Ethereum contract.
type AtomicSwapCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AtomicSwapTransactor is an auto generated write-only Go binding around an Ethereum contract.
type AtomicSwapTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AtomicSwapFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type AtomicSwapFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AtomicSwapSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type AtomicSwapSession struct {
	Contract     *AtomicSwap       // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// AtomicSwapCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type AtomicSwapCallerSession struct {
	Contract *AtomicSwapCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts     // Call options to use throughout this session
}

// AtomicSwapTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type AtomicSwapTransactorSession struct {
	Contract     *AtomicSwapTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// AtomicSwapRaw is an auto generated low-level Go binding around an Ethereum contract.
type AtomicSwapRaw struct {
	Contract *AtomicSwap // Generic contract binding to access the raw methods on
}

// AtomicSwapCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type AtomicSwapCallerRaw struct {
	Contract *AtomicSwapCaller // Generic read-only contract binding to access the raw methods on
}

// AtomicSwapTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type AtomicSwapTransactorRaw struct {
	Contract *AtomicSwapTransactor // Generic write-only contract binding to access the raw methods on
}

// NewAtomicSwap creates a new instance of AtomicSwap, bound to a specific deployed contract.
func NewAtomicSwap(address common.Address, backend bind.ContractBackend) (*AtomicSwap, error) {
	contract, err := bindAtomicSwap(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &AtomicSwap{AtomicSwapCaller: AtomicSwapCaller{contract: contract}, AtomicSwapTransactor: AtomicSwapTransactor{contract: contract}, AtomicSwapFilterer: AtomicSwapFilterer{contract: contract}}, nil
}

// NewAtomicSwapCaller creates a new read-only instance of AtomicSwap, bound to a specific deployed contract.
func NewAtomicSwapCaller(address common.Address, caller bind.ContractCaller) (*AtomicSwapCaller, error) {
	contract, err := bindAtomicSwap(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &AtomicSwapCaller{contract: contract}, nil
}

// NewAtomicSwapTransactor creates a new write-only instance of AtomicSwap, bound to a specific deployed contract.
func NewAtomicSwapTransactor(address common.Address, transactor bind.ContractTransactor) (*AtomicSwapTransactor, error) {
	contract, err := bindAtomicSwap(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &AtomicSwapTransactor{contract: contract}, nil
}

// NewAtomicSwapFilterer creates a new log filterer instance of AtomicSwap, bound to a specific deployed contract.
func NewAtomicSwapFilterer(address common.Address, filterer bind.ContractFilterer) (*AtomicSwapFilterer, error) {
	contract, err := bindAtomicSwap(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &AtomicSwapFilterer{contract: contract}, nil
}

// bindAtomicSwap binds a generic wrapper to an already deployed contract.
func bindAtomicSwap(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := AtomicSwapMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AtomicSwap *AtomicSwapRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AtomicSwap.Contract.AtomicSwapCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AtomicSwap *AtomicSwapRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AtomicSwap.Contract.AtomicSwapTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AtomicSwap *AtomicSwapRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AtomicSwap.Contract.AtomicSwapTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AtomicSwap *AtomicSwapCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AtomicSwap.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AtomicSwap *AtomicSwapTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AtomicSwap.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AtomicSwap *AtomicSwapTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AtomicSwap.Contract.contract.Transact(opts, method, params...)
}

// Redeem is a paid mutator transaction binding the contract method 0x3841185e.
//
// Solidity: function redeem(address _token, bytes _secret) returns()
func (_AtomicSwap *AtomicSwapTransactor) Redeem(opts *bind.TransactOpts, _token common.Address, _secret []byte) (*types.Transaction, error) {
	return _AtomicSwap.contract.Transact(opts, "redeem", _token, _secret)
}

// Redeem is a paid mutator transaction binding the contract method 0x3841185e.
//
// Solidity: function redeem(address _token, bytes _secret) returns()
func (_AtomicSwap *AtomicSwapSession) Redeem(_token common.Address, _secret []byte) (*types.Transaction, error) {
	return _AtomicSwap.Contract.Redeem(&_AtomicSwap.TransactOpts, _token, _secret)
}

// Redeem is a paid mutator transaction binding the contract method 0x3841185e.
//
// Solidity: function redeem(address _token, bytes _secret) returns()
func (_AtomicSwap *AtomicSwapTransactorSession) Redeem(_token common.Address, _secret []byte) (*types.Transaction, error) {
	return _AtomicSwap.Contract.Redeem(&_AtomicSwap.TransactOpts, _token, _secret)
}

// Refund is a paid mutator transaction binding the contract method 0xfa89401a.
//
// Solidity: function refund(address _token) returns()
func (_AtomicSwap *AtomicSwapTransactor) Refund(opts *bind.TransactOpts, _token common.Address) (*types.Transaction, error) {
	return _AtomicSwap.contract.Transact(opts, "refund", _token)
}

// Refund is a paid mutator transaction binding the contract method 0xfa89401a.
//
// Solidity: function refund(address _token) returns()
func (_AtomicSwap *AtomicSwapSession) Refund(_token common.Address) (*types.Transaction, error) {
	return _AtomicSwap.Contract.Refund(&_AtomicSwap.TransactOpts, _token)
}

// Refund is a paid mutator transaction binding the contract method 0xfa89401a.
//
// Solidity: function refund(address _token) returns()
func (_AtomicSwap *AtomicSwapTransactorSession) Refund(_token common.Address) (*types.Transaction, error) {
	return _AtomicSwap.Contract.Refund(&_AtomicSwap.TransactOpts, _token)
}

// AtomicSwapRedeemedIterator is returned from FilterRedeemed and is used to iterate over the raw logs and unpacked data for Redeemed events raised by the AtomicSwap contract.
type AtomicSwapRedeemedIterator struct {
	Event *AtomicSwapRedeemed // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *AtomicSwapRedeemedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AtomicSwapRedeemed)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(AtomicSwapRedeemed)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *AtomicSwapRedeemedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AtomicSwapRedeemedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AtomicSwapRedeemed represents a Redeemed event raised by the AtomicSwap contract.
type AtomicSwapRedeemed struct {
	Secret string
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterRedeemed is a free log retrieval operation binding the contract event 0x434fa9d5634e23e27249731c499be7ac2854b2717ff296496d611e3157f03642.
//
// Solidity: event Redeemed(string _secret)
func (_AtomicSwap *AtomicSwapFilterer) FilterRedeemed(opts *bind.FilterOpts) (*AtomicSwapRedeemedIterator, error) {

	logs, sub, err := _AtomicSwap.contract.FilterLogs(opts, "Redeemed")
	if err != nil {
		return nil, err
	}
	return &AtomicSwapRedeemedIterator{contract: _AtomicSwap.contract, event: "Redeemed", logs: logs, sub: sub}, nil
}

// WatchRedeemed is a free log subscription operation binding the contract event 0x434fa9d5634e23e27249731c499be7ac2854b2717ff296496d611e3157f03642.
//
// Solidity: event Redeemed(string _secret)
func (_AtomicSwap *AtomicSwapFilterer) WatchRedeemed(opts *bind.WatchOpts, sink chan<- *AtomicSwapRedeemed) (event.Subscription, error) {

	logs, sub, err := _AtomicSwap.contract.WatchLogs(opts, "Redeemed")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AtomicSwapRedeemed)
				if err := _AtomicSwap.contract.UnpackLog(event, "Redeemed", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseRedeemed is a log parse operation binding the contract event 0x434fa9d5634e23e27249731c499be7ac2854b2717ff296496d611e3157f03642.
//
// Solidity: event Redeemed(string _secret)
func (_AtomicSwap *AtomicSwapFilterer) ParseRedeemed(log types.Log) (*AtomicSwapRedeemed, error) {
	event := new(AtomicSwapRedeemed)
	if err := _AtomicSwap.contract.UnpackLog(event, "Redeemed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
