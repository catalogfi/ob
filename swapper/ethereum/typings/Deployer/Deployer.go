// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package Deployer

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

// DeployerMetaData contains all meta data concerning the Deployer contract.
var DeployerMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_feeCollector\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newFeeCollector\",\"type\":\"address\"}],\"name\":\"FeeCollectorUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Paused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Unpaused\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"tokenAddr\",\"type\":\"address\"}],\"name\":\"addToken\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"redeemer\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"refunder\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"tokenAddr\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"secretHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"expiry\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"computeAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"redeemer\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"refunder\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"tokenAddr\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"secretHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"expiry\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"deploy\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"redeemer\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"refunder\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"tokenAddr\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"secretHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"expiry\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"secret\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"deployAndRedeem\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"paused\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"tokenAddr\",\"type\":\"address\"}],\"name\":\"removeToken\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"unpause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newFeeCollector\",\"type\":\"address\"}],\"name\":\"updateFeeCollector\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b5060405161193f38038061193f83398101604081905261002f916100ba565b6100383361006a565b6000805460ff60a01b19169055600180546001600160a01b0319166001600160a01b03929092169190911790556100ea565b600080546001600160a01b038381166001600160a01b0319831681178455604051919092169283917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e09190a35050565b6000602082840312156100cc57600080fd5b81516001600160a01b03811681146100e357600080fd5b9392505050565b611846806100f96000396000f3fe608060405234801561001057600080fd5b50600436106100b45760003560e01c80638456cb59116100715780638456cb59146101265780638da5cb5b1461012e578063d2c35ce814610153578063d48bfca714610166578063e200eed414610179578063f2fde38b1461018c57600080fd5b806305fa6427146100b95780633f4ba83a146100ce578063443ae517146100d65780635c975abb146100e95780635fa7b5841461010b578063715018a61461011e575b600080fd5b6100cc6100c7366004610af1565b61019f565b005b6100cc610372565b6100cc6100e4366004610bf1565b610384565b600054600160a01b900460ff1660405190151581526020015b60405180910390f35b6100cc610119366004610c50565b6104ef565b6100cc610518565b6100cc61052a565b6000546001600160a01b03165b6040516001600160a01b039091168152602001610102565b6100cc610161366004610c50565b61053a565b6100cc610174366004610c50565b610596565b61013b610187366004610bf1565b6105c2565b6100cc61019a366004610c50565b610718565b8686864386116101ca5760405162461bcd60e51b81526004016101c190610c6b565b60405180910390fd5b6000604051806020016101dc90610ab2565b601f1982820381018352601f909101166040819052600154610216918e918e918e916001600160a01b0316908e908e908d90602001610cba565b60408051601f19818403018152908290526102349291602001610d22565b604051602081830303815290604052905061025160008983610791565b6001600160a01b0316639945e3d3876040518263ffffffff1660e01b815260040161027c9190610d51565b600060405180830381600087803b15801561029657600080fd5b505af11580156102aa573d6000803e3d6000fd5b5050506001600160a01b03851691506102d790505760405162461bcd60e51b81526004016101c190610d84565b6001600160a01b0383166102fd5760405162461bcd60e51b81526004016101c190610dcd565b826001600160a01b0316826001600160a01b03160361032e5760405162461bcd60e51b81526004016101c190610e17565b6001600160a01b03811660009081526002602052604090205460ff166103665760405162461bcd60e51b81526004016101c190610e66565b50505050505050505050565b61037a61089c565b6103826108f6565b565b8585854385116103a65760405162461bcd60e51b81526004016101c190610c6b565b6000604051806020016103b890610ab2565b601f1982820381018352601f9091011660408190526001546103f2918d918d918d916001600160a01b0316908d908d908d90602001610cba565b60408051601f19818403018152908290526104109291602001610d22565b604051602081830303815290604052905061042d60008883610791565b50506001600160a01b0383166104555760405162461bcd60e51b81526004016101c190610d84565b6001600160a01b03831661047b5760405162461bcd60e51b81526004016101c190610dcd565b826001600160a01b0316826001600160a01b0316036104ac5760405162461bcd60e51b81526004016101c190610e17565b6001600160a01b03811660009081526002602052604090205460ff166104e45760405162461bcd60e51b81526004016101c190610e66565b505050505050505050565b6104f761089c565b6001600160a01b03166000908152600260205260409020805460ff19169055565b61052061089c565b610382600061094b565b61053261089c565b61038261099b565b61054261089c565b600180546001600160a01b0319166001600160a01b0383169081179091556040519081527fe5693914d19c789bdee50a362998c0bc8d035a835f9871da5d51152f0582c34f9060200160405180910390a150565b61059e61089c565b6001600160a01b03166000908152600260205260409020805460ff19166001179055565b60008686866000604051806020016105d990610ab2565b601f1982820381018352601f909101166040819052600154610613918e918e918e916001600160a01b0316908e908e908e90602001610cba565b60408051601f19818403018152908290526106319291602001610d22565b60405160208183030381529060405290506106538882805190602001206109de565b9450506001600160a01b03831661067c5760405162461bcd60e51b81526004016101c190610d84565b6001600160a01b0383166106a25760405162461bcd60e51b81526004016101c190610dcd565b826001600160a01b0316826001600160a01b0316036106d35760405162461bcd60e51b81526004016101c190610e17565b6001600160a01b03811660009081526002602052604090205460ff1661070b5760405162461bcd60e51b81526004016101c190610e66565b5050509695505050505050565b61072061089c565b6001600160a01b0381166107855760405162461bcd60e51b815260206004820152602660248201527f4f776e61626c653a206e6577206f776e657220697320746865207a65726f206160448201526564647265737360d01b60648201526084016101c1565b61078e8161094b565b50565b6000834710156107e35760405162461bcd60e51b815260206004820152601d60248201527f437265617465323a20696e73756666696369656e742062616c616e636500000060448201526064016101c1565b81516000036108345760405162461bcd60e51b815260206004820181905260248201527f437265617465323a2062797465636f6465206c656e677468206973207a65726f60448201526064016101c1565b8282516020840186f590506001600160a01b0381166108955760405162461bcd60e51b815260206004820152601960248201527f437265617465323a204661696c6564206f6e206465706c6f790000000000000060448201526064016101c1565b9392505050565b6000546001600160a01b031633146103825760405162461bcd60e51b815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e657260448201526064016101c1565b6108fe6109eb565b6000805460ff60a01b191690557f5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa335b6040516001600160a01b03909116815260200160405180910390a1565b600080546001600160a01b038381166001600160a01b0319831681178455604051919092169283917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e09190a35050565b6109a3610a3b565b6000805460ff60a01b1916600160a01b1790557f62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a25861092e3390565b6000610895838330610a88565b600054600160a01b900460ff166103825760405162461bcd60e51b815260206004820152601460248201527314185d5cd8589b194e881b9bdd081c185d5cd95960621b60448201526064016101c1565b600054600160a01b900460ff16156103825760405162461bcd60e51b815260206004820152601060248201526f14185d5cd8589b194e881c185d5cd95960821b60448201526064016101c1565b6000604051836040820152846020820152828152600b8101905060ff815360559020949350505050565b61097380610e9e83390190565b80356001600160a01b0381168114610ad657600080fd5b919050565b634e487b7160e01b600052604160045260246000fd5b600080600080600080600060e0888a031215610b0c57600080fd5b610b1588610abf565b9650610b2360208901610abf565b9550610b3160408901610abf565b9450606088013593506080880135925060a088013567ffffffffffffffff80821115610b5c57600080fd5b818a0191508a601f830112610b7057600080fd5b813581811115610b8257610b82610adb565b604051601f8201601f19908116603f01168101908382118183101715610baa57610baa610adb565b816040528281528d6020848701011115610bc357600080fd5b82602086016020830137600060208483010152809650505050505060c0880135905092959891949750929550565b60008060008060008060c08789031215610c0a57600080fd5b610c1387610abf565b9550610c2160208801610abf565b9450610c2f60408801610abf565b9350606087013592506080870135915060a087013590509295509295509295565b600060208284031215610c6257600080fd5b61089582610abf565b6020808252602f908201527f4465706c6f7965723a20657870697279206d757374206265203e20637572726560408201526e373a10313637b1b590373ab6b132b960891b606082015260800190565b6001600160a01b03978816815295871660208701529386166040860152919094166060840152608083019390935260a082019290925260c081019190915260e00190565b60005b83811015610d19578181015183820152602001610d01565b50506000910152565b60008351610d34818460208801610cfe565b835190830190610d48818360208801610cfe565b01949350505050565b6020815260008251806020840152610d70816040850160208701610cfe565b601f01601f19169190910160400192915050565b60208082526029908201527f4465706c6f7965723a2072656465656d65722063616e6e6f74206265206e756c6040820152686c206164647265737360b81b606082015260800190565b6020808252602a908201527f4465706c6f7965723a20696e69746961746f722063616e6e6f74206265206e756040820152696c6c206164647265737360b01b606082015260800190565b6020808252602f908201527f4465706c6f7965723a20696e69746961746f722063616e6e6f7420626520657160408201526e3ab0b6103a37903932b232b2b6b2b960891b606082015260800190565b6020808252601d908201527f4465706c6f7965723a20546f6b656e206e6f7420737570706f7274656400000060408201526060019056fe61016060405234801561001157600080fd5b5060405161097338038061097383398101604081905261003091610080565b6001600160a01b0396871660805294861660a05292851660c052610100526101209190915290911660e052610140526100f1565b80516001600160a01b038116811461007b57600080fd5b919050565b600080600080600080600060e0888a03121561009b57600080fd5b6100a488610064565b96506100b260208901610064565b95506100c060408901610064565b94506100ce60608901610064565b93506080880151925060a0880151915060c0880151905092959891949750929550565b60805160a05160c05160e05161010051610120516101405161081c610157600039600081816101d301526102bc01526000605a015260006101070152600061026d01526000818161032f01526104de0152600060dd01526000610297015261081c6000f3fe608060405234801561001057600080fd5b50600436106100365760003560e01c8063590e1ae31461003b5780639945e3d314610045575b600080fd5b610043610058565b005b61004361005336600461068b565b610105565b7f000000000000000000000000000000000000000000000000000000000000000043116100cc5760405162461bcd60e51b815260206004820152601c60248201527f41746f6d6963537761703a206c6f636b206e6f7420657870697265640000000060448201526064015b60405180910390fd5b60006100d66102eb565b90506101027f00000000000000000000000000000000000000000000000000000000000000008261048e565b50565b7f0000000000000000000000000000000000000000000000000000000000000000600283836040516101389291906106fd565b602060405180830381855afa158015610155573d6000803e3d6000fd5b5050506040513d601f19601f82011682018060405250810190610178919061070d565b146101c55760405162461bcd60e51b815260206004820152601a60248201527f41746f6d6963537761703a2073656372657420696e76616c696400000000000060448201526064016100c3565b60006101cf6102eb565b90507f000000000000000000000000000000000000000000000000000000000000000081101561024c5760405162461bcd60e51b815260206004820152602260248201527f41746f6d6963537761703a20636f6e7472616374206e6f7420696e6974696174604482015261195960f21b60648201526084016100c3565b600061271061025c601e8461073c565b6102669190610759565b90506102927f00000000000000000000000000000000000000000000000000000000000000008261048e565b6102e57f00000000000000000000000000000000000000000000000000000000000000006102e0837f000000000000000000000000000000000000000000000000000000000000000061077b565b61048e565b50505050565b604080513060248083019190915282518083039091018152604490910182526020810180516001600160e01b03166370a0823160e01b1790529051600091829182917f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031691610362919061078e565b6000604051808303816000865af19150503d806000811461039f576040519150601f19603f3d011682016040523d82523d6000602084013e6103a4565b606091505b50915091508161040a5760405162461bcd60e51b815260206004820152602b60248201527f41746f6d6963537761703a2045524332302062616c616e63654f66206469642060448201526a1b9bdd081cdd58d8d9595960aa1b60648201526084016100c3565b60008151116104735760405162461bcd60e51b815260206004820152602f60248201527f41746f6d6963537761703a2045524332302062616c616e63654f66206469642060448201526e6e6f742072657475726e206461746160881b60648201526084016100c3565b80806020019051810190610487919061070d565b9250505090565b604080516001600160a01b038481166024830152604480830185905283518084039091018152606490920183526020820180516001600160e01b031663a9059cbb60e01b179052915160009283927f00000000000000000000000000000000000000000000000000000000000000009091169161050b919061078e565b6000604051808303816000865af19150503d8060008114610548576040519150601f19603f3d011682016040523d82523d6000602084013e61054d565b606091505b5091509150816105b95760405162461bcd60e51b815260206004820152603160248201527f41746f6d6963537761703a204552433230207472616e7366657220646964206e6044820152706f7420737563636565642028626f6f6c2960781b60648201526084016100c3565b60008151116106215760405162461bcd60e51b815260206004820152602e60248201527f41746f6d6963537761703a204552433230207472616e7366657220646964206e60448201526d6f742072657475726e206461746160901b60648201526084016100c3565b8080602001905181019061063591906107bd565b6102e55760405162461bcd60e51b815260206004820152602160248201527f41746f6d6963537761703a204552433230207472616e73666572206661696c656044820152601960fa1b60648201526084016100c3565b6000806020838503121561069e57600080fd5b823567ffffffffffffffff808211156106b657600080fd5b818501915085601f8301126106ca57600080fd5b8135818111156106d957600080fd5b8660208285010111156106eb57600080fd5b60209290920196919550909350505050565b8183823760009101908152919050565b60006020828403121561071f57600080fd5b5051919050565b634e487b7160e01b600052601160045260246000fd5b808202811582820484141761075357610753610726565b92915050565b60008261077657634e487b7160e01b600052601260045260246000fd5b500490565b8181038181111561075357610753610726565b6000825160005b818110156107af5760208186018101518583015201610795565b506000920191825250919050565b6000602082840312156107cf57600080fd5b815180151581146107df57600080fd5b939250505056fea2646970667358221220f9b83b7a9660f3abfbf581b10996600f0149769dc65854a56a5157a49bf5307664736f6c63430008120033a2646970667358221220c91655a4315075fc6625af5024e5752c80336a9580947444ee9c24b887879ba164736f6c63430008120033",
}

// DeployerABI is the input ABI used to generate the binding from.
// Deprecated: Use DeployerMetaData.ABI instead.
var DeployerABI = DeployerMetaData.ABI

// DeployerBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use DeployerMetaData.Bin instead.
var DeployerBin = DeployerMetaData.Bin

// DeployDeployer deploys a new Ethereum contract, binding an instance of Deployer to it.
func DeployDeployer(auth *bind.TransactOpts, backend bind.ContractBackend, _feeCollector common.Address) (common.Address, *types.Transaction, *Deployer, error) {
	parsed, err := DeployerMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(DeployerBin), backend, _feeCollector)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Deployer{DeployerCaller: DeployerCaller{contract: contract}, DeployerTransactor: DeployerTransactor{contract: contract}, DeployerFilterer: DeployerFilterer{contract: contract}}, nil
}

// Deployer is an auto generated Go binding around an Ethereum contract.
type Deployer struct {
	DeployerCaller     // Read-only binding to the contract
	DeployerTransactor // Write-only binding to the contract
	DeployerFilterer   // Log filterer for contract events
}

// DeployerCaller is an auto generated read-only Go binding around an Ethereum contract.
type DeployerCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DeployerTransactor is an auto generated write-only Go binding around an Ethereum contract.
type DeployerTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DeployerFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type DeployerFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DeployerSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type DeployerSession struct {
	Contract     *Deployer         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// DeployerCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type DeployerCallerSession struct {
	Contract *DeployerCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// DeployerTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type DeployerTransactorSession struct {
	Contract     *DeployerTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// DeployerRaw is an auto generated low-level Go binding around an Ethereum contract.
type DeployerRaw struct {
	Contract *Deployer // Generic contract binding to access the raw methods on
}

// DeployerCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type DeployerCallerRaw struct {
	Contract *DeployerCaller // Generic read-only contract binding to access the raw methods on
}

// DeployerTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type DeployerTransactorRaw struct {
	Contract *DeployerTransactor // Generic write-only contract binding to access the raw methods on
}

// NewDeployer creates a new instance of Deployer, bound to a specific deployed contract.
func NewDeployer(address common.Address, backend bind.ContractBackend) (*Deployer, error) {
	contract, err := bindDeployer(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Deployer{DeployerCaller: DeployerCaller{contract: contract}, DeployerTransactor: DeployerTransactor{contract: contract}, DeployerFilterer: DeployerFilterer{contract: contract}}, nil
}

// NewDeployerCaller creates a new read-only instance of Deployer, bound to a specific deployed contract.
func NewDeployerCaller(address common.Address, caller bind.ContractCaller) (*DeployerCaller, error) {
	contract, err := bindDeployer(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &DeployerCaller{contract: contract}, nil
}

// NewDeployerTransactor creates a new write-only instance of Deployer, bound to a specific deployed contract.
func NewDeployerTransactor(address common.Address, transactor bind.ContractTransactor) (*DeployerTransactor, error) {
	contract, err := bindDeployer(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &DeployerTransactor{contract: contract}, nil
}

// NewDeployerFilterer creates a new log filterer instance of Deployer, bound to a specific deployed contract.
func NewDeployerFilterer(address common.Address, filterer bind.ContractFilterer) (*DeployerFilterer, error) {
	contract, err := bindDeployer(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &DeployerFilterer{contract: contract}, nil
}

// bindDeployer binds a generic wrapper to an already deployed contract.
func bindDeployer(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(DeployerABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Deployer *DeployerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Deployer.Contract.DeployerCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Deployer *DeployerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Deployer.Contract.DeployerTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Deployer *DeployerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Deployer.Contract.DeployerTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Deployer *DeployerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Deployer.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Deployer *DeployerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Deployer.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Deployer *DeployerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Deployer.Contract.contract.Transact(opts, method, params...)
}

// ComputeAddress is a free data retrieval call binding the contract method 0xe200eed4.
//
// Solidity: function computeAddress(address redeemer, address refunder, address tokenAddr, bytes32 secretHash, uint256 expiry, uint256 amount) view returns(address)
func (_Deployer *DeployerCaller) ComputeAddress(opts *bind.CallOpts, redeemer common.Address, refunder common.Address, tokenAddr common.Address, secretHash [32]byte, expiry *big.Int, amount *big.Int) (common.Address, error) {
	var out []interface{}
	err := _Deployer.contract.Call(opts, &out, "computeAddress", redeemer, refunder, tokenAddr, secretHash, expiry, amount)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// ComputeAddress is a free data retrieval call binding the contract method 0xe200eed4.
//
// Solidity: function computeAddress(address redeemer, address refunder, address tokenAddr, bytes32 secretHash, uint256 expiry, uint256 amount) view returns(address)
func (_Deployer *DeployerSession) ComputeAddress(redeemer common.Address, refunder common.Address, tokenAddr common.Address, secretHash [32]byte, expiry *big.Int, amount *big.Int) (common.Address, error) {
	return _Deployer.Contract.ComputeAddress(&_Deployer.CallOpts, redeemer, refunder, tokenAddr, secretHash, expiry, amount)
}

// ComputeAddress is a free data retrieval call binding the contract method 0xe200eed4.
//
// Solidity: function computeAddress(address redeemer, address refunder, address tokenAddr, bytes32 secretHash, uint256 expiry, uint256 amount) view returns(address)
func (_Deployer *DeployerCallerSession) ComputeAddress(redeemer common.Address, refunder common.Address, tokenAddr common.Address, secretHash [32]byte, expiry *big.Int, amount *big.Int) (common.Address, error) {
	return _Deployer.Contract.ComputeAddress(&_Deployer.CallOpts, redeemer, refunder, tokenAddr, secretHash, expiry, amount)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Deployer *DeployerCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Deployer.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Deployer *DeployerSession) Owner() (common.Address, error) {
	return _Deployer.Contract.Owner(&_Deployer.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Deployer *DeployerCallerSession) Owner() (common.Address, error) {
	return _Deployer.Contract.Owner(&_Deployer.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_Deployer *DeployerCaller) Paused(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _Deployer.contract.Call(opts, &out, "paused")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_Deployer *DeployerSession) Paused() (bool, error) {
	return _Deployer.Contract.Paused(&_Deployer.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_Deployer *DeployerCallerSession) Paused() (bool, error) {
	return _Deployer.Contract.Paused(&_Deployer.CallOpts)
}

// AddToken is a paid mutator transaction binding the contract method 0xd48bfca7.
//
// Solidity: function addToken(address tokenAddr) returns()
func (_Deployer *DeployerTransactor) AddToken(opts *bind.TransactOpts, tokenAddr common.Address) (*types.Transaction, error) {
	return _Deployer.contract.Transact(opts, "addToken", tokenAddr)
}

// AddToken is a paid mutator transaction binding the contract method 0xd48bfca7.
//
// Solidity: function addToken(address tokenAddr) returns()
func (_Deployer *DeployerSession) AddToken(tokenAddr common.Address) (*types.Transaction, error) {
	return _Deployer.Contract.AddToken(&_Deployer.TransactOpts, tokenAddr)
}

// AddToken is a paid mutator transaction binding the contract method 0xd48bfca7.
//
// Solidity: function addToken(address tokenAddr) returns()
func (_Deployer *DeployerTransactorSession) AddToken(tokenAddr common.Address) (*types.Transaction, error) {
	return _Deployer.Contract.AddToken(&_Deployer.TransactOpts, tokenAddr)
}

// Deploy is a paid mutator transaction binding the contract method 0x443ae517.
//
// Solidity: function deploy(address redeemer, address refunder, address tokenAddr, bytes32 secretHash, uint256 expiry, uint256 amount) returns()
func (_Deployer *DeployerTransactor) Deploy(opts *bind.TransactOpts, redeemer common.Address, refunder common.Address, tokenAddr common.Address, secretHash [32]byte, expiry *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _Deployer.contract.Transact(opts, "deploy", redeemer, refunder, tokenAddr, secretHash, expiry, amount)
}

// Deploy is a paid mutator transaction binding the contract method 0x443ae517.
//
// Solidity: function deploy(address redeemer, address refunder, address tokenAddr, bytes32 secretHash, uint256 expiry, uint256 amount) returns()
func (_Deployer *DeployerSession) Deploy(redeemer common.Address, refunder common.Address, tokenAddr common.Address, secretHash [32]byte, expiry *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _Deployer.Contract.Deploy(&_Deployer.TransactOpts, redeemer, refunder, tokenAddr, secretHash, expiry, amount)
}

// Deploy is a paid mutator transaction binding the contract method 0x443ae517.
//
// Solidity: function deploy(address redeemer, address refunder, address tokenAddr, bytes32 secretHash, uint256 expiry, uint256 amount) returns()
func (_Deployer *DeployerTransactorSession) Deploy(redeemer common.Address, refunder common.Address, tokenAddr common.Address, secretHash [32]byte, expiry *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _Deployer.Contract.Deploy(&_Deployer.TransactOpts, redeemer, refunder, tokenAddr, secretHash, expiry, amount)
}

// DeployAndRedeem is a paid mutator transaction binding the contract method 0x05fa6427.
//
// Solidity: function deployAndRedeem(address redeemer, address refunder, address tokenAddr, bytes32 secretHash, uint256 expiry, bytes secret, uint256 amount) returns()
func (_Deployer *DeployerTransactor) DeployAndRedeem(opts *bind.TransactOpts, redeemer common.Address, refunder common.Address, tokenAddr common.Address, secretHash [32]byte, expiry *big.Int, secret []byte, amount *big.Int) (*types.Transaction, error) {
	return _Deployer.contract.Transact(opts, "deployAndRedeem", redeemer, refunder, tokenAddr, secretHash, expiry, secret, amount)
}

// DeployAndRedeem is a paid mutator transaction binding the contract method 0x05fa6427.
//
// Solidity: function deployAndRedeem(address redeemer, address refunder, address tokenAddr, bytes32 secretHash, uint256 expiry, bytes secret, uint256 amount) returns()
func (_Deployer *DeployerSession) DeployAndRedeem(redeemer common.Address, refunder common.Address, tokenAddr common.Address, secretHash [32]byte, expiry *big.Int, secret []byte, amount *big.Int) (*types.Transaction, error) {
	return _Deployer.Contract.DeployAndRedeem(&_Deployer.TransactOpts, redeemer, refunder, tokenAddr, secretHash, expiry, secret, amount)
}

// DeployAndRedeem is a paid mutator transaction binding the contract method 0x05fa6427.
//
// Solidity: function deployAndRedeem(address redeemer, address refunder, address tokenAddr, bytes32 secretHash, uint256 expiry, bytes secret, uint256 amount) returns()
func (_Deployer *DeployerTransactorSession) DeployAndRedeem(redeemer common.Address, refunder common.Address, tokenAddr common.Address, secretHash [32]byte, expiry *big.Int, secret []byte, amount *big.Int) (*types.Transaction, error) {
	return _Deployer.Contract.DeployAndRedeem(&_Deployer.TransactOpts, redeemer, refunder, tokenAddr, secretHash, expiry, secret, amount)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_Deployer *DeployerTransactor) Pause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Deployer.contract.Transact(opts, "pause")
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_Deployer *DeployerSession) Pause() (*types.Transaction, error) {
	return _Deployer.Contract.Pause(&_Deployer.TransactOpts)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_Deployer *DeployerTransactorSession) Pause() (*types.Transaction, error) {
	return _Deployer.Contract.Pause(&_Deployer.TransactOpts)
}

// RemoveToken is a paid mutator transaction binding the contract method 0x5fa7b584.
//
// Solidity: function removeToken(address tokenAddr) returns()
func (_Deployer *DeployerTransactor) RemoveToken(opts *bind.TransactOpts, tokenAddr common.Address) (*types.Transaction, error) {
	return _Deployer.contract.Transact(opts, "removeToken", tokenAddr)
}

// RemoveToken is a paid mutator transaction binding the contract method 0x5fa7b584.
//
// Solidity: function removeToken(address tokenAddr) returns()
func (_Deployer *DeployerSession) RemoveToken(tokenAddr common.Address) (*types.Transaction, error) {
	return _Deployer.Contract.RemoveToken(&_Deployer.TransactOpts, tokenAddr)
}

// RemoveToken is a paid mutator transaction binding the contract method 0x5fa7b584.
//
// Solidity: function removeToken(address tokenAddr) returns()
func (_Deployer *DeployerTransactorSession) RemoveToken(tokenAddr common.Address) (*types.Transaction, error) {
	return _Deployer.Contract.RemoveToken(&_Deployer.TransactOpts, tokenAddr)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Deployer *DeployerTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Deployer.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Deployer *DeployerSession) RenounceOwnership() (*types.Transaction, error) {
	return _Deployer.Contract.RenounceOwnership(&_Deployer.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Deployer *DeployerTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _Deployer.Contract.RenounceOwnership(&_Deployer.TransactOpts)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Deployer *DeployerTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _Deployer.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Deployer *DeployerSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Deployer.Contract.TransferOwnership(&_Deployer.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Deployer *DeployerTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Deployer.Contract.TransferOwnership(&_Deployer.TransactOpts, newOwner)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_Deployer *DeployerTransactor) Unpause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Deployer.contract.Transact(opts, "unpause")
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_Deployer *DeployerSession) Unpause() (*types.Transaction, error) {
	return _Deployer.Contract.Unpause(&_Deployer.TransactOpts)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_Deployer *DeployerTransactorSession) Unpause() (*types.Transaction, error) {
	return _Deployer.Contract.Unpause(&_Deployer.TransactOpts)
}

// UpdateFeeCollector is a paid mutator transaction binding the contract method 0xd2c35ce8.
//
// Solidity: function updateFeeCollector(address newFeeCollector) returns()
func (_Deployer *DeployerTransactor) UpdateFeeCollector(opts *bind.TransactOpts, newFeeCollector common.Address) (*types.Transaction, error) {
	return _Deployer.contract.Transact(opts, "updateFeeCollector", newFeeCollector)
}

// UpdateFeeCollector is a paid mutator transaction binding the contract method 0xd2c35ce8.
//
// Solidity: function updateFeeCollector(address newFeeCollector) returns()
func (_Deployer *DeployerSession) UpdateFeeCollector(newFeeCollector common.Address) (*types.Transaction, error) {
	return _Deployer.Contract.UpdateFeeCollector(&_Deployer.TransactOpts, newFeeCollector)
}

// UpdateFeeCollector is a paid mutator transaction binding the contract method 0xd2c35ce8.
//
// Solidity: function updateFeeCollector(address newFeeCollector) returns()
func (_Deployer *DeployerTransactorSession) UpdateFeeCollector(newFeeCollector common.Address) (*types.Transaction, error) {
	return _Deployer.Contract.UpdateFeeCollector(&_Deployer.TransactOpts, newFeeCollector)
}

// DeployerFeeCollectorUpdatedIterator is returned from FilterFeeCollectorUpdated and is used to iterate over the raw logs and unpacked data for FeeCollectorUpdated events raised by the Deployer contract.
type DeployerFeeCollectorUpdatedIterator struct {
	Event *DeployerFeeCollectorUpdated // Event containing the contract specifics and raw log

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
func (it *DeployerFeeCollectorUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DeployerFeeCollectorUpdated)
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
		it.Event = new(DeployerFeeCollectorUpdated)
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
func (it *DeployerFeeCollectorUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DeployerFeeCollectorUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DeployerFeeCollectorUpdated represents a FeeCollectorUpdated event raised by the Deployer contract.
type DeployerFeeCollectorUpdated struct {
	NewFeeCollector common.Address
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterFeeCollectorUpdated is a free log retrieval operation binding the contract event 0xe5693914d19c789bdee50a362998c0bc8d035a835f9871da5d51152f0582c34f.
//
// Solidity: event FeeCollectorUpdated(address newFeeCollector)
func (_Deployer *DeployerFilterer) FilterFeeCollectorUpdated(opts *bind.FilterOpts) (*DeployerFeeCollectorUpdatedIterator, error) {

	logs, sub, err := _Deployer.contract.FilterLogs(opts, "FeeCollectorUpdated")
	if err != nil {
		return nil, err
	}
	return &DeployerFeeCollectorUpdatedIterator{contract: _Deployer.contract, event: "FeeCollectorUpdated", logs: logs, sub: sub}, nil
}

// WatchFeeCollectorUpdated is a free log subscription operation binding the contract event 0xe5693914d19c789bdee50a362998c0bc8d035a835f9871da5d51152f0582c34f.
//
// Solidity: event FeeCollectorUpdated(address newFeeCollector)
func (_Deployer *DeployerFilterer) WatchFeeCollectorUpdated(opts *bind.WatchOpts, sink chan<- *DeployerFeeCollectorUpdated) (event.Subscription, error) {

	logs, sub, err := _Deployer.contract.WatchLogs(opts, "FeeCollectorUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DeployerFeeCollectorUpdated)
				if err := _Deployer.contract.UnpackLog(event, "FeeCollectorUpdated", log); err != nil {
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

// ParseFeeCollectorUpdated is a log parse operation binding the contract event 0xe5693914d19c789bdee50a362998c0bc8d035a835f9871da5d51152f0582c34f.
//
// Solidity: event FeeCollectorUpdated(address newFeeCollector)
func (_Deployer *DeployerFilterer) ParseFeeCollectorUpdated(log types.Log) (*DeployerFeeCollectorUpdated, error) {
	event := new(DeployerFeeCollectorUpdated)
	if err := _Deployer.contract.UnpackLog(event, "FeeCollectorUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DeployerOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the Deployer contract.
type DeployerOwnershipTransferredIterator struct {
	Event *DeployerOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *DeployerOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DeployerOwnershipTransferred)
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
		it.Event = new(DeployerOwnershipTransferred)
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
func (it *DeployerOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DeployerOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DeployerOwnershipTransferred represents a OwnershipTransferred event raised by the Deployer contract.
type DeployerOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Deployer *DeployerFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*DeployerOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Deployer.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &DeployerOwnershipTransferredIterator{contract: _Deployer.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Deployer *DeployerFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *DeployerOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Deployer.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DeployerOwnershipTransferred)
				if err := _Deployer.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Deployer *DeployerFilterer) ParseOwnershipTransferred(log types.Log) (*DeployerOwnershipTransferred, error) {
	event := new(DeployerOwnershipTransferred)
	if err := _Deployer.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DeployerPausedIterator is returned from FilterPaused and is used to iterate over the raw logs and unpacked data for Paused events raised by the Deployer contract.
type DeployerPausedIterator struct {
	Event *DeployerPaused // Event containing the contract specifics and raw log

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
func (it *DeployerPausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DeployerPaused)
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
		it.Event = new(DeployerPaused)
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
func (it *DeployerPausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DeployerPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DeployerPaused represents a Paused event raised by the Deployer contract.
type DeployerPaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterPaused is a free log retrieval operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_Deployer *DeployerFilterer) FilterPaused(opts *bind.FilterOpts) (*DeployerPausedIterator, error) {

	logs, sub, err := _Deployer.contract.FilterLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return &DeployerPausedIterator{contract: _Deployer.contract, event: "Paused", logs: logs, sub: sub}, nil
}

// WatchPaused is a free log subscription operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_Deployer *DeployerFilterer) WatchPaused(opts *bind.WatchOpts, sink chan<- *DeployerPaused) (event.Subscription, error) {

	logs, sub, err := _Deployer.contract.WatchLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DeployerPaused)
				if err := _Deployer.contract.UnpackLog(event, "Paused", log); err != nil {
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

// ParsePaused is a log parse operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_Deployer *DeployerFilterer) ParsePaused(log types.Log) (*DeployerPaused, error) {
	event := new(DeployerPaused)
	if err := _Deployer.contract.UnpackLog(event, "Paused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DeployerUnpausedIterator is returned from FilterUnpaused and is used to iterate over the raw logs and unpacked data for Unpaused events raised by the Deployer contract.
type DeployerUnpausedIterator struct {
	Event *DeployerUnpaused // Event containing the contract specifics and raw log

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
func (it *DeployerUnpausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DeployerUnpaused)
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
		it.Event = new(DeployerUnpaused)
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
func (it *DeployerUnpausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DeployerUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DeployerUnpaused represents a Unpaused event raised by the Deployer contract.
type DeployerUnpaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterUnpaused is a free log retrieval operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_Deployer *DeployerFilterer) FilterUnpaused(opts *bind.FilterOpts) (*DeployerUnpausedIterator, error) {

	logs, sub, err := _Deployer.contract.FilterLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return &DeployerUnpausedIterator{contract: _Deployer.contract, event: "Unpaused", logs: logs, sub: sub}, nil
}

// WatchUnpaused is a free log subscription operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_Deployer *DeployerFilterer) WatchUnpaused(opts *bind.WatchOpts, sink chan<- *DeployerUnpaused) (event.Subscription, error) {

	logs, sub, err := _Deployer.contract.WatchLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DeployerUnpaused)
				if err := _Deployer.contract.UnpackLog(event, "Unpaused", log); err != nil {
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

// ParseUnpaused is a log parse operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_Deployer *DeployerFilterer) ParseUnpaused(log types.Log) (*DeployerUnpaused, error) {
	event := new(DeployerUnpaused)
	if err := _Deployer.contract.UnpackLog(event, "Unpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
