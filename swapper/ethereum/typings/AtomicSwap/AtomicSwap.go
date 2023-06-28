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
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"secrectHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Initiated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"secrectHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"_secret\",\"type\":\"bytes\"}],\"name\":\"Redeemed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"secrectHash\",\"type\":\"bytes32\"}],\"name\":\"Refunded\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_redeemer\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_expiry\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"_secretHash\",\"type\":\"bytes32\"}],\"name\":\"initiate\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"_secret\",\"type\":\"bytes\"}],\"name\":\"redeem\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_secretHash\",\"type\":\"bytes32\"}],\"name\":\"refund\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60a060405234801561001057600080fd5b50604051610a2d380380610a2d83398101604081905261002f91610040565b6001600160a01b0316608052610070565b60006020828403121561005257600080fd5b81516001600160a01b038116811461006957600080fd5b9392505050565b608051610994610099600039600081816101a40152818161036f015261071901526109946000f3fe608060405234801561001057600080fd5b50600436106100415760003560e01c80637249fbb61461004657806397ffc7ae1461005b5780639945e3d31461006e575b600080fd5b6100596100543660046107c7565b610081565b005b6100596100693660046107e0565b610241565b61005961007c366004610827565b6105e0565b600081815260208190526040902080546001600160a01b03166100eb5760405162461bcd60e51b815260206004820152601e60248201527f41746f6d6963537761703a206f72646572206e6f7420696e697461746564000060448201526064015b60405180910390fd5b600481015460ff16156101105760405162461bcd60e51b81526004016100e290610899565b806002015443116101635760405162461bcd60e51b815260206004820152601c60248201527f41746f6d6963537761703a206c6f636b206e6f7420657870697265640000000060448201526064016100e2565b6004818101805460ff19166001908117909155820154600383015460405163a9059cbb60e01b81526001600160a01b039283169381019390935260248301527f0000000000000000000000000000000000000000000000000000000000000000169063a9059cbb906044016020604051808303816000875af11580156101ed573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061021191906108dd565b5060405182907ffe509803c09416b28ff3d8f690c8b0c61462a892c46d5430c8fb20abe472daf090600090a25050565b60008181526020818152604091829020825160a08101845281546001600160a01b039081168252600183015416928101929092526002810154928201929092526003820154606082015260049091015460ff161580156080830152859133918691906102ef5760405162461bcd60e51b815260206004820152601f60248201527f41746f6d6963537761703a2063616e6e6f74207265757365207365637265740060448201526064016100e2565b80516001600160a01b0316156103175760405162461bcd60e51b81526004016100e290610899565b6040805160a0810182526001600160a01b038a8116825233602083018190528284018b9052606083018a90526000608084015292516323b872dd60e01b815260048101939093523060248401526044830189905290917f0000000000000000000000000000000000000000000000000000000000000000909116906323b872dd906064016020604051808303816000875af11580156103ba573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906103de91906108dd565b5060008681526020818152604091829020835181546001600160a01b03199081166001600160a01b0392831617835585840151600184018054909216921691909117905583830151600282015560608401516003820181905560808501516004909201805460ff191692151592909217909155915191825287917fbd7231421af354010a8dc99d32bc090722c773f05c06893cafffbdc19d9b5a89910160405180910390a250506001600160a01b0383166104ef5760405162461bcd60e51b815260206004820152602b60248201527f41746f6d6963537761703a2072656465656d65722063616e6e6f74206265206e60448201526a756c6c206164647265737360a81b60648201526084016100e2565b826001600160a01b0316826001600160a01b03160361056a5760405162461bcd60e51b815260206004820152603160248201527f41746f6d6963537761703a20696e69746961746f722063616e6e6f742062652060448201527032b8bab0b6103a37903932b232b2b6b2b960791b60648201526084016100e2565b4381116105d75760405162461bcd60e51b815260206004820152603560248201527f41746f6d6963537761703a206578706972792063616e6e6f74206265206c6f776044820152746572207468616e2063757272656e7420626c6f636b60581b60648201526084016100e2565b50505050505050565b6000600283836040516105f4929190610906565b602060405180830381855afa158015610611573d6000803e3d6000fd5b5050506040513d601f19601f820116820180604052508101906106349190610916565b60008181526020819052604090208054919250906001600160a01b03166106b75760405162461bcd60e51b815260206004820152603160248201527f41746f6d6963537761703a20696e76616c696420736563726574206f72206f7260448201527019195c881b9bdd081a5b9a5d1a585d1959607a1b60648201526084016100e2565b600481015460ff16156106dc5760405162461bcd60e51b81526004016100e290610899565b6004818101805460ff191660011790558154600383015460405163a9059cbb60e01b81526001600160a01b039283169381019390935260248301527f0000000000000000000000000000000000000000000000000000000000000000169063a9059cbb906044016020604051808303816000875af1158015610762573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061078691906108dd565b50817f866c33f43c7dda3105124ae616b2a42ff25811f48048edbb4ab215c59563b1e685856040516107b992919061092f565b60405180910390a250505050565b6000602082840312156107d957600080fd5b5035919050565b600080600080608085870312156107f657600080fd5b84356001600160a01b038116811461080d57600080fd5b966020860135965060408601359560600135945092505050565b6000806020838503121561083a57600080fd5b823567ffffffffffffffff8082111561085257600080fd5b818501915085601f83011261086657600080fd5b81358181111561087557600080fd5b86602082850101111561088757600080fd5b60209290920196919550909350505050565b60208082526024908201527f41746f6d6963537761703a206f7264657220616c72656164792066756c6c66696040820152631b1b195960e21b606082015260800190565b6000602082840312156108ef57600080fd5b815180151581146108ff57600080fd5b9392505050565b8183823760009101908152919050565b60006020828403121561092857600080fd5b5051919050565b60208152816020820152818360408301376000818301604090810191909152601f909201601f1916010191905056fea26469706673582212207d2aec99649175d521949c31478e23470cb5f7b6bd51854c4a5746cd1005b15a64736f6c63430008120033",
}

// AtomicSwapABI is the input ABI used to generate the binding from.
// Deprecated: Use AtomicSwapMetaData.ABI instead.
var AtomicSwapABI = AtomicSwapMetaData.ABI

// AtomicSwapBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use AtomicSwapMetaData.Bin instead.
var AtomicSwapBin = AtomicSwapMetaData.Bin

// DeployAtomicSwap deploys a new Ethereum contract, binding an instance of AtomicSwap to it.
func DeployAtomicSwap(auth *bind.TransactOpts, backend bind.ContractBackend, _token common.Address) (common.Address, *types.Transaction, *AtomicSwap, error) {
	parsed, err := AtomicSwapMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(AtomicSwapBin), backend, _token)
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

// Initiate is a paid mutator transaction binding the contract method 0x97ffc7ae.
//
// Solidity: function initiate(address _redeemer, uint256 _expiry, uint256 _amount, bytes32 _secretHash) returns()
func (_AtomicSwap *AtomicSwapTransactor) Initiate(opts *bind.TransactOpts, _redeemer common.Address, _expiry *big.Int, _amount *big.Int, _secretHash [32]byte) (*types.Transaction, error) {
	return _AtomicSwap.contract.Transact(opts, "initiate", _redeemer, _expiry, _amount, _secretHash)
}

// Initiate is a paid mutator transaction binding the contract method 0x97ffc7ae.
//
// Solidity: function initiate(address _redeemer, uint256 _expiry, uint256 _amount, bytes32 _secretHash) returns()
func (_AtomicSwap *AtomicSwapSession) Initiate(_redeemer common.Address, _expiry *big.Int, _amount *big.Int, _secretHash [32]byte) (*types.Transaction, error) {
	return _AtomicSwap.Contract.Initiate(&_AtomicSwap.TransactOpts, _redeemer, _expiry, _amount, _secretHash)
}

// Initiate is a paid mutator transaction binding the contract method 0x97ffc7ae.
//
// Solidity: function initiate(address _redeemer, uint256 _expiry, uint256 _amount, bytes32 _secretHash) returns()
func (_AtomicSwap *AtomicSwapTransactorSession) Initiate(_redeemer common.Address, _expiry *big.Int, _amount *big.Int, _secretHash [32]byte) (*types.Transaction, error) {
	return _AtomicSwap.Contract.Initiate(&_AtomicSwap.TransactOpts, _redeemer, _expiry, _amount, _secretHash)
}

// Redeem is a paid mutator transaction binding the contract method 0x9945e3d3.
//
// Solidity: function redeem(bytes _secret) returns()
func (_AtomicSwap *AtomicSwapTransactor) Redeem(opts *bind.TransactOpts, _secret []byte) (*types.Transaction, error) {
	return _AtomicSwap.contract.Transact(opts, "redeem", _secret)
}

// Redeem is a paid mutator transaction binding the contract method 0x9945e3d3.
//
// Solidity: function redeem(bytes _secret) returns()
func (_AtomicSwap *AtomicSwapSession) Redeem(_secret []byte) (*types.Transaction, error) {
	return _AtomicSwap.Contract.Redeem(&_AtomicSwap.TransactOpts, _secret)
}

// Redeem is a paid mutator transaction binding the contract method 0x9945e3d3.
//
// Solidity: function redeem(bytes _secret) returns()
func (_AtomicSwap *AtomicSwapTransactorSession) Redeem(_secret []byte) (*types.Transaction, error) {
	return _AtomicSwap.Contract.Redeem(&_AtomicSwap.TransactOpts, _secret)
}

// Refund is a paid mutator transaction binding the contract method 0x7249fbb6.
//
// Solidity: function refund(bytes32 _secretHash) returns()
func (_AtomicSwap *AtomicSwapTransactor) Refund(opts *bind.TransactOpts, _secretHash [32]byte) (*types.Transaction, error) {
	return _AtomicSwap.contract.Transact(opts, "refund", _secretHash)
}

// Refund is a paid mutator transaction binding the contract method 0x7249fbb6.
//
// Solidity: function refund(bytes32 _secretHash) returns()
func (_AtomicSwap *AtomicSwapSession) Refund(_secretHash [32]byte) (*types.Transaction, error) {
	return _AtomicSwap.Contract.Refund(&_AtomicSwap.TransactOpts, _secretHash)
}

// Refund is a paid mutator transaction binding the contract method 0x7249fbb6.
//
// Solidity: function refund(bytes32 _secretHash) returns()
func (_AtomicSwap *AtomicSwapTransactorSession) Refund(_secretHash [32]byte) (*types.Transaction, error) {
	return _AtomicSwap.Contract.Refund(&_AtomicSwap.TransactOpts, _secretHash)
}

// AtomicSwapInitiatedIterator is returned from FilterInitiated and is used to iterate over the raw logs and unpacked data for Initiated events raised by the AtomicSwap contract.
type AtomicSwapInitiatedIterator struct {
	Event *AtomicSwapInitiated // Event containing the contract specifics and raw log

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
func (it *AtomicSwapInitiatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AtomicSwapInitiated)
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
		it.Event = new(AtomicSwapInitiated)
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
func (it *AtomicSwapInitiatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AtomicSwapInitiatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AtomicSwapInitiated represents a Initiated event raised by the AtomicSwap contract.
type AtomicSwapInitiated struct {
	SecrectHash [32]byte
	Amount      *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterInitiated is a free log retrieval operation binding the contract event 0xbd7231421af354010a8dc99d32bc090722c773f05c06893cafffbdc19d9b5a89.
//
// Solidity: event Initiated(bytes32 indexed secrectHash, uint256 amount)
func (_AtomicSwap *AtomicSwapFilterer) FilterInitiated(opts *bind.FilterOpts, secrectHash [][32]byte) (*AtomicSwapInitiatedIterator, error) {

	var secrectHashRule []interface{}
	for _, secrectHashItem := range secrectHash {
		secrectHashRule = append(secrectHashRule, secrectHashItem)
	}

	logs, sub, err := _AtomicSwap.contract.FilterLogs(opts, "Initiated", secrectHashRule)
	if err != nil {
		return nil, err
	}
	return &AtomicSwapInitiatedIterator{contract: _AtomicSwap.contract, event: "Initiated", logs: logs, sub: sub}, nil
}

// WatchInitiated is a free log subscription operation binding the contract event 0xbd7231421af354010a8dc99d32bc090722c773f05c06893cafffbdc19d9b5a89.
//
// Solidity: event Initiated(bytes32 indexed secrectHash, uint256 amount)
func (_AtomicSwap *AtomicSwapFilterer) WatchInitiated(opts *bind.WatchOpts, sink chan<- *AtomicSwapInitiated, secrectHash [][32]byte) (event.Subscription, error) {

	var secrectHashRule []interface{}
	for _, secrectHashItem := range secrectHash {
		secrectHashRule = append(secrectHashRule, secrectHashItem)
	}

	logs, sub, err := _AtomicSwap.contract.WatchLogs(opts, "Initiated", secrectHashRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AtomicSwapInitiated)
				if err := _AtomicSwap.contract.UnpackLog(event, "Initiated", log); err != nil {
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

// ParseInitiated is a log parse operation binding the contract event 0xbd7231421af354010a8dc99d32bc090722c773f05c06893cafffbdc19d9b5a89.
//
// Solidity: event Initiated(bytes32 indexed secrectHash, uint256 amount)
func (_AtomicSwap *AtomicSwapFilterer) ParseInitiated(log types.Log) (*AtomicSwapInitiated, error) {
	event := new(AtomicSwapInitiated)
	if err := _AtomicSwap.contract.UnpackLog(event, "Initiated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
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
	SecrectHash [32]byte
	Secret      []byte
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterRedeemed is a free log retrieval operation binding the contract event 0x866c33f43c7dda3105124ae616b2a42ff25811f48048edbb4ab215c59563b1e6.
//
// Solidity: event Redeemed(bytes32 indexed secrectHash, bytes _secret)
func (_AtomicSwap *AtomicSwapFilterer) FilterRedeemed(opts *bind.FilterOpts, secrectHash [][32]byte) (*AtomicSwapRedeemedIterator, error) {

	var secrectHashRule []interface{}
	for _, secrectHashItem := range secrectHash {
		secrectHashRule = append(secrectHashRule, secrectHashItem)
	}

	logs, sub, err := _AtomicSwap.contract.FilterLogs(opts, "Redeemed", secrectHashRule)
	if err != nil {
		return nil, err
	}
	return &AtomicSwapRedeemedIterator{contract: _AtomicSwap.contract, event: "Redeemed", logs: logs, sub: sub}, nil
}

// WatchRedeemed is a free log subscription operation binding the contract event 0x866c33f43c7dda3105124ae616b2a42ff25811f48048edbb4ab215c59563b1e6.
//
// Solidity: event Redeemed(bytes32 indexed secrectHash, bytes _secret)
func (_AtomicSwap *AtomicSwapFilterer) WatchRedeemed(opts *bind.WatchOpts, sink chan<- *AtomicSwapRedeemed, secrectHash [][32]byte) (event.Subscription, error) {

	var secrectHashRule []interface{}
	for _, secrectHashItem := range secrectHash {
		secrectHashRule = append(secrectHashRule, secrectHashItem)
	}

	logs, sub, err := _AtomicSwap.contract.WatchLogs(opts, "Redeemed", secrectHashRule)
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

// ParseRedeemed is a log parse operation binding the contract event 0x866c33f43c7dda3105124ae616b2a42ff25811f48048edbb4ab215c59563b1e6.
//
// Solidity: event Redeemed(bytes32 indexed secrectHash, bytes _secret)
func (_AtomicSwap *AtomicSwapFilterer) ParseRedeemed(log types.Log) (*AtomicSwapRedeemed, error) {
	event := new(AtomicSwapRedeemed)
	if err := _AtomicSwap.contract.UnpackLog(event, "Redeemed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AtomicSwapRefundedIterator is returned from FilterRefunded and is used to iterate over the raw logs and unpacked data for Refunded events raised by the AtomicSwap contract.
type AtomicSwapRefundedIterator struct {
	Event *AtomicSwapRefunded // Event containing the contract specifics and raw log

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
func (it *AtomicSwapRefundedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AtomicSwapRefunded)
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
		it.Event = new(AtomicSwapRefunded)
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
func (it *AtomicSwapRefundedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AtomicSwapRefundedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AtomicSwapRefunded represents a Refunded event raised by the AtomicSwap contract.
type AtomicSwapRefunded struct {
	SecrectHash [32]byte
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterRefunded is a free log retrieval operation binding the contract event 0xfe509803c09416b28ff3d8f690c8b0c61462a892c46d5430c8fb20abe472daf0.
//
// Solidity: event Refunded(bytes32 indexed secrectHash)
func (_AtomicSwap *AtomicSwapFilterer) FilterRefunded(opts *bind.FilterOpts, secrectHash [][32]byte) (*AtomicSwapRefundedIterator, error) {

	var secrectHashRule []interface{}
	for _, secrectHashItem := range secrectHash {
		secrectHashRule = append(secrectHashRule, secrectHashItem)
	}

	logs, sub, err := _AtomicSwap.contract.FilterLogs(opts, "Refunded", secrectHashRule)
	if err != nil {
		return nil, err
	}
	return &AtomicSwapRefundedIterator{contract: _AtomicSwap.contract, event: "Refunded", logs: logs, sub: sub}, nil
}

// WatchRefunded is a free log subscription operation binding the contract event 0xfe509803c09416b28ff3d8f690c8b0c61462a892c46d5430c8fb20abe472daf0.
//
// Solidity: event Refunded(bytes32 indexed secrectHash)
func (_AtomicSwap *AtomicSwapFilterer) WatchRefunded(opts *bind.WatchOpts, sink chan<- *AtomicSwapRefunded, secrectHash [][32]byte) (event.Subscription, error) {

	var secrectHashRule []interface{}
	for _, secrectHashItem := range secrectHash {
		secrectHashRule = append(secrectHashRule, secrectHashItem)
	}

	logs, sub, err := _AtomicSwap.contract.WatchLogs(opts, "Refunded", secrectHashRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AtomicSwapRefunded)
				if err := _AtomicSwap.contract.UnpackLog(event, "Refunded", log); err != nil {
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

// ParseRefunded is a log parse operation binding the contract event 0xfe509803c09416b28ff3d8f690c8b0c61462a892c46d5430c8fb20abe472daf0.
//
// Solidity: event Refunded(bytes32 indexed secrectHash)
func (_AtomicSwap *AtomicSwapFilterer) ParseRefunded(log types.Log) (*AtomicSwapRefunded, error) {
	event := new(AtomicSwapRefunded)
	if err := _AtomicSwap.contract.UnpackLog(event, "Refunded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
