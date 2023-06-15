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
)

// AtomicSwapMetaData contains all meta data concerning the AtomicSwap contract.
var AtomicSwapMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_redeemer\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_refunder\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"_secretHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"_expiry\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"_secret\",\"type\":\"bytes\"}],\"name\":\"redeemed\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"_secret\",\"type\":\"bytes\"}],\"name\":\"execute\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x6101006040523480156200001257600080fd5b5060405162000c6a38038062000c6a83398181016040528101906200003891906200019a565b8373ffffffffffffffffffffffffffffffffffffffff1660808173ffffffffffffffffffffffffffffffffffffffff16815250508273ffffffffffffffffffffffffffffffffffffffff1660a08173ffffffffffffffffffffffffffffffffffffffff16815250508160c081815250508060e08181525050505050506200020c565b600080fd5b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b6000620000ec82620000bf565b9050919050565b620000fe81620000df565b81146200010a57600080fd5b50565b6000815190506200011e81620000f3565b92915050565b6000819050919050565b620001398162000124565b81146200014557600080fd5b50565b60008151905062000159816200012e565b92915050565b6000819050919050565b62000174816200015f565b81146200018057600080fd5b50565b600081519050620001948162000169565b92915050565b60008060008060808587031215620001b757620001b6620000ba565b5b6000620001c7878288016200010d565b9450506020620001da878288016200010d565b9350506040620001ed8782880162000148565b9250506060620002008782880162000183565b91505092959194509250565b60805160a05160c05160e051610a2762000243600039600060e201526000604e0152600061010e0152600060800152610a276000f3fe608060405234801561001057600080fd5b506004361061002b5760003560e01c80631cff79cd14610030575b600080fd5b61004a600480360381019061004591906105d6565b61004c565b005b7f00000000000000000000000000000000000000000000000000000000000000008180519060200120036100e0576100a4827f0000000000000000000000000000000000000000000000000000000000000000610177565b7f803bbe225132e3ac86fc914e09473a7f35093547f3c3466472eb38b6ea7b7e9d816040516100d391906106b1565b60405180910390a1610173565b7f000000000000000000000000000000000000000000000000000000000000000043111561013757610132827f0000000000000000000000000000000000000000000000000000000000000000610177565b610172565b6040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161016990610730565b60405180910390fd5b5b5050565b6000808373ffffffffffffffffffffffffffffffffffffffff166370a08231306040516024016101a7919061075f565b6040516020818303038152906040529060e01b6020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff83818316178352505050506040516101f591906107b6565b6000604051808303816000865af19150503d8060008114610232576040519150601f19603f3d011682016040523d82523d6000602084013e610237565b606091505b50915091508161027c576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016102739061083f565b60405180910390fd5b6000818060200190518101906102929190610895565b90506000811115610405576000808673ffffffffffffffffffffffffffffffffffffffff1663a9059cbb87856040516024016102cf9291906108d1565b6040516020818303038152906040529060e01b6020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff838183161783525050505060405161031d91906107b6565b6000604051808303816000865af19150503d806000811461035a576040519150601f19603f3d011682016040523d82523d6000602084013e61035f565b606091505b5091509150816103a4576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161039b9061096c565b60405180910390fd5b60008151111561040257808060200190518101906103c291906109c4565b610401576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016103f89061096c565b60405180910390fd5b5b50505b8373ffffffffffffffffffffffffffffffffffffffff16ff5b6000604051905090565b600080fd5b600080fd5b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b600061045d82610432565b9050919050565b61046d81610452565b811461047857600080fd5b50565b60008135905061048a81610464565b92915050565b600080fd5b600080fd5b6000601f19601f8301169050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6104e38261049a565b810181811067ffffffffffffffff82111715610502576105016104ab565b5b80604052505050565b600061051561041e565b905061052182826104da565b919050565b600067ffffffffffffffff821115610541576105406104ab565b5b61054a8261049a565b9050602081019050919050565b82818337600083830152505050565b600061057961057484610526565b61050b565b90508281526020810184848401111561059557610594610495565b5b6105a0848285610557565b509392505050565b600082601f8301126105bd576105bc610490565b5b81356105cd848260208601610566565b91505092915050565b600080604083850312156105ed576105ec610428565b5b60006105fb8582860161047b565b925050602083013567ffffffffffffffff81111561061c5761061b61042d565b5b610628858286016105a8565b9150509250929050565b600081519050919050565b600082825260208201905092915050565b60005b8381101561066c578082015181840152602081019050610651565b60008484015250505050565b600061068382610632565b61068d818561063d565b935061069d81856020860161064e565b6106a68161049a565b840191505092915050565b600060208201905081810360008301526106cb8184610678565b905092915050565b600082825260208201905092915050565b7f696e76616c696420736563726574206f72206578706972790000000000000000600082015250565b600061071a6018836106d3565b9150610725826106e4565b602082019050919050565b600060208201905081810360008301526107498161070d565b9050919050565b61075981610452565b82525050565b60006020820190506107746000830184610750565b92915050565b600081905092915050565b600061079082610632565b61079a818561077a565b93506107aa81856020860161064e565b80840191505092915050565b60006107c28284610785565b915081905092915050565b7f41746f6d6963537761703a2045524332302062616c616e63654f66206469642060008201527f6e6f742073756363656564000000000000000000000000000000000000000000602082015250565b6000610829602b836106d3565b9150610834826107cd565b604082019050919050565b600060208201905081810360008301526108588161081c565b9050919050565b6000819050919050565b6108728161085f565b811461087d57600080fd5b50565b60008151905061088f81610869565b92915050565b6000602082840312156108ab576108aa610428565b5b60006108b984828501610880565b91505092915050565b6108cb8161085f565b82525050565b60006040820190506108e66000830185610750565b6108f360208301846108c2565b9392505050565b7f41746f6d6963537761703a204552433230207472616e7366657220646964206e60008201527f6f74207375636365656400000000000000000000000000000000000000000000602082015250565b6000610956602a836106d3565b9150610961826108fa565b604082019050919050565b6000602082019050818103600083015261098581610949565b9050919050565b60008115159050919050565b6109a18161098c565b81146109ac57600080fd5b50565b6000815190506109be81610998565b92915050565b6000602082840312156109da576109d9610428565b5b60006109e8848285016109af565b9150509291505056fea26469706673582212207b9f6cabd842e224fa7e87f5c7ea9c0c70f2e5f9844b6dbd792e667dc9a4148664736f6c63430008110033",
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
	parsed, err := abi.JSON(strings.NewReader(AtomicSwapABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
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

// Execute is a paid mutator transaction binding the contract method 0x1cff79cd.
//
// Solidity: function execute(address _token, bytes _secret) returns()
func (_AtomicSwap *AtomicSwapTransactor) Execute(opts *bind.TransactOpts, _token common.Address, _secret []byte) (*types.Transaction, error) {
	return _AtomicSwap.contract.Transact(opts, "execute", _token, _secret)
}

// Execute is a paid mutator transaction binding the contract method 0x1cff79cd.
//
// Solidity: function execute(address _token, bytes _secret) returns()
func (_AtomicSwap *AtomicSwapSession) Execute(_token common.Address, _secret []byte) (*types.Transaction, error) {
	return _AtomicSwap.Contract.Execute(&_AtomicSwap.TransactOpts, _token, _secret)
}

// Execute is a paid mutator transaction binding the contract method 0x1cff79cd.
//
// Solidity: function execute(address _token, bytes _secret) returns()
func (_AtomicSwap *AtomicSwapTransactorSession) Execute(_token common.Address, _secret []byte) (*types.Transaction, error) {
	return _AtomicSwap.Contract.Execute(&_AtomicSwap.TransactOpts, _token, _secret)
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
	Secret []byte
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterRedeemed is a free log retrieval operation binding the contract event 0x803bbe225132e3ac86fc914e09473a7f35093547f3c3466472eb38b6ea7b7e9d.
//
// Solidity: event redeemed(bytes _secret)
func (_AtomicSwap *AtomicSwapFilterer) FilterRedeemed(opts *bind.FilterOpts) (*AtomicSwapRedeemedIterator, error) {

	logs, sub, err := _AtomicSwap.contract.FilterLogs(opts, "redeemed")
	if err != nil {
		return nil, err
	}
	return &AtomicSwapRedeemedIterator{contract: _AtomicSwap.contract, event: "redeemed", logs: logs, sub: sub}, nil
}

// WatchRedeemed is a free log subscription operation binding the contract event 0x803bbe225132e3ac86fc914e09473a7f35093547f3c3466472eb38b6ea7b7e9d.
//
// Solidity: event redeemed(bytes _secret)
func (_AtomicSwap *AtomicSwapFilterer) WatchRedeemed(opts *bind.WatchOpts, sink chan<- *AtomicSwapRedeemed) (event.Subscription, error) {

	logs, sub, err := _AtomicSwap.contract.WatchLogs(opts, "redeemed")
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
				if err := _AtomicSwap.contract.UnpackLog(event, "redeemed", log); err != nil {
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

// ParseRedeemed is a log parse operation binding the contract event 0x803bbe225132e3ac86fc914e09473a7f35093547f3c3466472eb38b6ea7b7e9d.
//
// Solidity: event redeemed(bytes _secret)
func (_AtomicSwap *AtomicSwapFilterer) ParseRedeemed(log types.Log) (*AtomicSwapRedeemed, error) {
	event := new(AtomicSwapRedeemed)
	if err := _AtomicSwap.contract.UnpackLog(event, "redeemed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
