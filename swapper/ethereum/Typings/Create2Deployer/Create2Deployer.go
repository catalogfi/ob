// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package Create2Deployer

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

// Create2DeployerMetaData contains all meta data concerning the Create2Deployer contract.
var Create2DeployerMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Paused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Unpaused\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"salt\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"codeHash\",\"type\":\"bytes32\"}],\"name\":\"computeAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"salt\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"codeHash\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"deployer\",\"type\":\"address\"}],\"name\":\"computeAddressWithDeployer\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"salt\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"code\",\"type\":\"bytes\"}],\"name\":\"deploy\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"salt\",\"type\":\"bytes32\"}],\"name\":\"deployERC1820Implementer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"addresspayable\",\"name\":\"payoutAddress\",\"type\":\"address\"}],\"name\":\"killCreate2Deployer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"paused\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"unpause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
	Bin: "0x608060405234801561001057600080fd5b5061002d61002261004c60201b60201c565b61005460201b60201c565b60008060146101000a81548160ff021916908315150217905550610118565b600033905090565b60008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff169050816000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508173ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a35050565b6115a0806101276000396000f3fe6080604052600436106100a05760003560e01c80636447045411610064578063644704541461019157806366cfa057146101ba578063715018a6146101e35780638456cb59146101fa5780638da5cb5b14610211578063f2fde38b1461023c576100a7565b8063076c37b2146100ac5780633f4ba83a146100d5578063481286e6146100ec57806356299481146101295780635c975abb14610166576100a7565b366100a757005b600080fd5b3480156100b857600080fd5b506100d360048036038101906100ce9190610b0e565b610265565b005b3480156100e157600080fd5b506100ea6102e0565b005b3480156100f857600080fd5b50610113600480360381019061010e9190610b4e565b610366565b6040516101209190610bcf565b60405180910390f35b34801561013557600080fd5b50610150600480360381019061014b9190610c16565b61037a565b60405161015d9190610bcf565b60405180910390f35b34801561017257600080fd5b5061017b610390565b6040516101889190610c84565b60405180910390f35b34801561019d57600080fd5b506101b860048036038101906101b39190610cdd565b6103a6565b005b3480156101c657600080fd5b506101e160048036038101906101dc9190610e50565b610482565b005b3480156101ef57600080fd5b506101f86104db565b005b34801561020657600080fd5b5061020f610563565b005b34801561021d57600080fd5b506102266105e9565b6040516102339190610bcf565b60405180910390f35b34801561024857600080fd5b50610263600480360381019061025e9190610ebf565b610612565b005b61026d610390565b156102ad576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016102a490610f49565b60405180910390fd5b6102db8282604051806020016102c290610a81565b6020820181038252601f19601f82011660405250610709565b505050565b6102e8610818565b73ffffffffffffffffffffffffffffffffffffffff166103066105e9565b73ffffffffffffffffffffffffffffffffffffffff161461035c576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161035390610fb5565b60405180910390fd5b610364610820565b565b600061037283836108c1565b905092915050565b60006103878484846108d6565b90509392505050565b60008060149054906101000a900460ff16905090565b6103ae610818565b73ffffffffffffffffffffffffffffffffffffffff166103cc6105e9565b73ffffffffffffffffffffffffffffffffffffffff1614610422576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161041990610fb5565b60405180910390fd5b8073ffffffffffffffffffffffffffffffffffffffff166108fc479081150290604051600060405180830381858888f19350505050158015610468573d6000803e3d6000fd5b508073ffffffffffffffffffffffffffffffffffffffff16ff5b61048a610390565b156104ca576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016104c190610f49565b60405180910390fd5b6104d5838383610709565b50505050565b6104e3610818565b73ffffffffffffffffffffffffffffffffffffffff166105016105e9565b73ffffffffffffffffffffffffffffffffffffffff1614610557576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161054e90610fb5565b60405180910390fd5b610561600061091a565b565b61056b610818565b73ffffffffffffffffffffffffffffffffffffffff166105896105e9565b73ffffffffffffffffffffffffffffffffffffffff16146105df576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016105d690610fb5565b60405180910390fd5b6105e76109de565b565b60008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905090565b61061a610818565b73ffffffffffffffffffffffffffffffffffffffff166106386105e9565b73ffffffffffffffffffffffffffffffffffffffff161461068e576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161068590610fb5565b60405180910390fd5b600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff16036106fd576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016106f490611047565b60405180910390fd5b6107068161091a565b50565b6000808447101561074f576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610746906110b3565b60405180910390fd5b6000835103610793576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161078a9061111f565b60405180910390fd5b8383516020850187f59050600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff160361080d576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016108049061118b565b60405180910390fd5b809150509392505050565b600033905090565b610828610390565b610867576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161085e906111f7565b60405180910390fd5b60008060146101000a81548160ff0219169083151502179055507f5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa6108aa610818565b6040516108b79190610bcf565b60405180910390a1565b60006108ce8383306108d6565b905092915050565b60008060ff60f81b8386866040516020016108f494939291906112cd565b6040516020818303038152906040528051906020012090508060001c9150509392505050565b60008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff169050816000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508173ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a35050565b6109e6610390565b15610a26576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610a1d90610f49565b60405180910390fd5b6001600060146101000a81548160ff0219169083151502179055507f62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258610a6a610818565b604051610a779190610bcf565b60405180910390a1565b61024f8061131c83390190565b6000604051905090565b600080fd5b600080fd5b6000819050919050565b610ab581610aa2565b8114610ac057600080fd5b50565b600081359050610ad281610aac565b92915050565b6000819050919050565b610aeb81610ad8565b8114610af657600080fd5b50565b600081359050610b0881610ae2565b92915050565b60008060408385031215610b2557610b24610a98565b5b6000610b3385828601610ac3565b9250506020610b4485828601610af9565b9150509250929050565b60008060408385031215610b6557610b64610a98565b5b6000610b7385828601610af9565b9250506020610b8485828601610af9565b9150509250929050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b6000610bb982610b8e565b9050919050565b610bc981610bae565b82525050565b6000602082019050610be46000830184610bc0565b92915050565b610bf381610bae565b8114610bfe57600080fd5b50565b600081359050610c1081610bea565b92915050565b600080600060608486031215610c2f57610c2e610a98565b5b6000610c3d86828701610af9565b9350506020610c4e86828701610af9565b9250506040610c5f86828701610c01565b9150509250925092565b60008115159050919050565b610c7e81610c69565b82525050565b6000602082019050610c996000830184610c75565b92915050565b6000610caa82610b8e565b9050919050565b610cba81610c9f565b8114610cc557600080fd5b50565b600081359050610cd781610cb1565b92915050565b600060208284031215610cf357610cf2610a98565b5b6000610d0184828501610cc8565b91505092915050565b600080fd5b600080fd5b6000601f19601f8301169050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b610d5d82610d14565b810181811067ffffffffffffffff82111715610d7c57610d7b610d25565b5b80604052505050565b6000610d8f610a8e565b9050610d9b8282610d54565b919050565b600067ffffffffffffffff821115610dbb57610dba610d25565b5b610dc482610d14565b9050602081019050919050565b82818337600083830152505050565b6000610df3610dee84610da0565b610d85565b905082815260208101848484011115610e0f57610e0e610d0f565b5b610e1a848285610dd1565b509392505050565b600082601f830112610e3757610e36610d0a565b5b8135610e47848260208601610de0565b91505092915050565b600080600060608486031215610e6957610e68610a98565b5b6000610e7786828701610ac3565b9350506020610e8886828701610af9565b925050604084013567ffffffffffffffff811115610ea957610ea8610a9d565b5b610eb586828701610e22565b9150509250925092565b600060208284031215610ed557610ed4610a98565b5b6000610ee384828501610c01565b91505092915050565b600082825260208201905092915050565b7f5061757361626c653a2070617573656400000000000000000000000000000000600082015250565b6000610f33601083610eec565b9150610f3e82610efd565b602082019050919050565b60006020820190508181036000830152610f6281610f26565b9050919050565b7f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e6572600082015250565b6000610f9f602083610eec565b9150610faa82610f69565b602082019050919050565b60006020820190508181036000830152610fce81610f92565b9050919050565b7f4f776e61626c653a206e6577206f776e657220697320746865207a65726f206160008201527f6464726573730000000000000000000000000000000000000000000000000000602082015250565b6000611031602683610eec565b915061103c82610fd5565b604082019050919050565b6000602082019050818103600083015261106081611024565b9050919050565b7f437265617465323a20696e73756666696369656e742062616c616e6365000000600082015250565b600061109d601d83610eec565b91506110a882611067565b602082019050919050565b600060208201905081810360008301526110cc81611090565b9050919050565b7f437265617465323a2062797465636f6465206c656e677468206973207a65726f600082015250565b6000611109602083610eec565b9150611114826110d3565b602082019050919050565b60006020820190508181036000830152611138816110fc565b9050919050565b7f437265617465323a204661696c6564206f6e206465706c6f7900000000000000600082015250565b6000611175601983610eec565b91506111808261113f565b602082019050919050565b600060208201905081810360008301526111a481611168565b9050919050565b7f5061757361626c653a206e6f7420706175736564000000000000000000000000600082015250565b60006111e1601483610eec565b91506111ec826111ab565b602082019050919050565b60006020820190508181036000830152611210816111d4565b9050919050565b60007fff0000000000000000000000000000000000000000000000000000000000000082169050919050565b6000819050919050565b61125e61125982611217565b611243565b82525050565b60008160601b9050919050565b600061127c82611264565b9050919050565b600061128e82611271565b9050919050565b6112a66112a182610bae565b611283565b82525050565b6000819050919050565b6112c76112c282610ad8565b6112ac565b82525050565b60006112d9828761124d565b6001820191506112e98286611295565b6014820191506112f982856112b6565b60208201915061130982846112b6565b6020820191508190509594505050505056fe608060405234801561001057600080fd5b5061022f806100206000396000f3fe608060405234801561001057600080fd5b506004361061002b5760003560e01c8063249cb3fa14610030575b600080fd5b61004a6004803603810190610045919061018f565b610060565b60405161005791906101de565b60405180910390f35b600080600084815260200190815260200160002060008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff166100cc576000801b6100ee565b7fa2ef4600d742022d532d4747cb3547474667d6f13804902513b2ec01c848f4b45b905092915050565b600080fd5b6000819050919050565b61010e816100fb565b811461011957600080fd5b50565b60008135905061012b81610105565b92915050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b600061015c82610131565b9050919050565b61016c81610151565b811461017757600080fd5b50565b60008135905061018981610163565b92915050565b600080604083850312156101a6576101a56100f6565b5b60006101b48582860161011c565b92505060206101c58582860161017a565b9150509250929050565b6101d8816100fb565b82525050565b60006020820190506101f360008301846101cf565b9291505056fea2646970667358221220dec1c8f4fd26ac778a02644a0854437c7c7a8b00cc57f37abddd831ae41712ab64736f6c63430008110033a2646970667358221220bc59f893efbfcce9dcf5358435d61130e15eeeb76d723c75128289520d31cb1a64736f6c63430008110033",
}

// Create2DeployerABI is the input ABI used to generate the binding from.
// Deprecated: Use Create2DeployerMetaData.ABI instead.
var Create2DeployerABI = Create2DeployerMetaData.ABI

// Create2DeployerBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use Create2DeployerMetaData.Bin instead.
var Create2DeployerBin = Create2DeployerMetaData.Bin

// DeployCreate2Deployer deploys a new Ethereum contract, binding an instance of Create2Deployer to it.
func DeployCreate2Deployer(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Create2Deployer, error) {
	parsed, err := Create2DeployerMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(Create2DeployerBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Create2Deployer{Create2DeployerCaller: Create2DeployerCaller{contract: contract}, Create2DeployerTransactor: Create2DeployerTransactor{contract: contract}, Create2DeployerFilterer: Create2DeployerFilterer{contract: contract}}, nil
}

// Create2Deployer is an auto generated Go binding around an Ethereum contract.
type Create2Deployer struct {
	Create2DeployerCaller     // Read-only binding to the contract
	Create2DeployerTransactor // Write-only binding to the contract
	Create2DeployerFilterer   // Log filterer for contract events
}

// Create2DeployerCaller is an auto generated read-only Go binding around an Ethereum contract.
type Create2DeployerCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// Create2DeployerTransactor is an auto generated write-only Go binding around an Ethereum contract.
type Create2DeployerTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// Create2DeployerFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type Create2DeployerFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// Create2DeployerSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type Create2DeployerSession struct {
	Contract     *Create2Deployer  // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// Create2DeployerCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type Create2DeployerCallerSession struct {
	Contract *Create2DeployerCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts          // Call options to use throughout this session
}

// Create2DeployerTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type Create2DeployerTransactorSession struct {
	Contract     *Create2DeployerTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts          // Transaction auth options to use throughout this session
}

// Create2DeployerRaw is an auto generated low-level Go binding around an Ethereum contract.
type Create2DeployerRaw struct {
	Contract *Create2Deployer // Generic contract binding to access the raw methods on
}

// Create2DeployerCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type Create2DeployerCallerRaw struct {
	Contract *Create2DeployerCaller // Generic read-only contract binding to access the raw methods on
}

// Create2DeployerTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type Create2DeployerTransactorRaw struct {
	Contract *Create2DeployerTransactor // Generic write-only contract binding to access the raw methods on
}

// NewCreate2Deployer creates a new instance of Create2Deployer, bound to a specific deployed contract.
func NewCreate2Deployer(address common.Address, backend bind.ContractBackend) (*Create2Deployer, error) {
	contract, err := bindCreate2Deployer(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Create2Deployer{Create2DeployerCaller: Create2DeployerCaller{contract: contract}, Create2DeployerTransactor: Create2DeployerTransactor{contract: contract}, Create2DeployerFilterer: Create2DeployerFilterer{contract: contract}}, nil
}

// NewCreate2DeployerCaller creates a new read-only instance of Create2Deployer, bound to a specific deployed contract.
func NewCreate2DeployerCaller(address common.Address, caller bind.ContractCaller) (*Create2DeployerCaller, error) {
	contract, err := bindCreate2Deployer(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &Create2DeployerCaller{contract: contract}, nil
}

// NewCreate2DeployerTransactor creates a new write-only instance of Create2Deployer, bound to a specific deployed contract.
func NewCreate2DeployerTransactor(address common.Address, transactor bind.ContractTransactor) (*Create2DeployerTransactor, error) {
	contract, err := bindCreate2Deployer(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &Create2DeployerTransactor{contract: contract}, nil
}

// NewCreate2DeployerFilterer creates a new log filterer instance of Create2Deployer, bound to a specific deployed contract.
func NewCreate2DeployerFilterer(address common.Address, filterer bind.ContractFilterer) (*Create2DeployerFilterer, error) {
	contract, err := bindCreate2Deployer(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &Create2DeployerFilterer{contract: contract}, nil
}

// bindCreate2Deployer binds a generic wrapper to an already deployed contract.
func bindCreate2Deployer(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(Create2DeployerABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Create2Deployer *Create2DeployerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Create2Deployer.Contract.Create2DeployerCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Create2Deployer *Create2DeployerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Create2Deployer.Contract.Create2DeployerTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Create2Deployer *Create2DeployerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Create2Deployer.Contract.Create2DeployerTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Create2Deployer *Create2DeployerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Create2Deployer.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Create2Deployer *Create2DeployerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Create2Deployer.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Create2Deployer *Create2DeployerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Create2Deployer.Contract.contract.Transact(opts, method, params...)
}

// ComputeAddress is a free data retrieval call binding the contract method 0x481286e6.
//
// Solidity: function computeAddress(bytes32 salt, bytes32 codeHash) view returns(address)
func (_Create2Deployer *Create2DeployerCaller) ComputeAddress(opts *bind.CallOpts, salt [32]byte, codeHash [32]byte) (common.Address, error) {
	var out []interface{}
	err := _Create2Deployer.contract.Call(opts, &out, "computeAddress", salt, codeHash)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// ComputeAddress is a free data retrieval call binding the contract method 0x481286e6.
//
// Solidity: function computeAddress(bytes32 salt, bytes32 codeHash) view returns(address)
func (_Create2Deployer *Create2DeployerSession) ComputeAddress(salt [32]byte, codeHash [32]byte) (common.Address, error) {
	return _Create2Deployer.Contract.ComputeAddress(&_Create2Deployer.CallOpts, salt, codeHash)
}

// ComputeAddress is a free data retrieval call binding the contract method 0x481286e6.
//
// Solidity: function computeAddress(bytes32 salt, bytes32 codeHash) view returns(address)
func (_Create2Deployer *Create2DeployerCallerSession) ComputeAddress(salt [32]byte, codeHash [32]byte) (common.Address, error) {
	return _Create2Deployer.Contract.ComputeAddress(&_Create2Deployer.CallOpts, salt, codeHash)
}

// ComputeAddressWithDeployer is a free data retrieval call binding the contract method 0x56299481.
//
// Solidity: function computeAddressWithDeployer(bytes32 salt, bytes32 codeHash, address deployer) pure returns(address)
func (_Create2Deployer *Create2DeployerCaller) ComputeAddressWithDeployer(opts *bind.CallOpts, salt [32]byte, codeHash [32]byte, deployer common.Address) (common.Address, error) {
	var out []interface{}
	err := _Create2Deployer.contract.Call(opts, &out, "computeAddressWithDeployer", salt, codeHash, deployer)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// ComputeAddressWithDeployer is a free data retrieval call binding the contract method 0x56299481.
//
// Solidity: function computeAddressWithDeployer(bytes32 salt, bytes32 codeHash, address deployer) pure returns(address)
func (_Create2Deployer *Create2DeployerSession) ComputeAddressWithDeployer(salt [32]byte, codeHash [32]byte, deployer common.Address) (common.Address, error) {
	return _Create2Deployer.Contract.ComputeAddressWithDeployer(&_Create2Deployer.CallOpts, salt, codeHash, deployer)
}

// ComputeAddressWithDeployer is a free data retrieval call binding the contract method 0x56299481.
//
// Solidity: function computeAddressWithDeployer(bytes32 salt, bytes32 codeHash, address deployer) pure returns(address)
func (_Create2Deployer *Create2DeployerCallerSession) ComputeAddressWithDeployer(salt [32]byte, codeHash [32]byte, deployer common.Address) (common.Address, error) {
	return _Create2Deployer.Contract.ComputeAddressWithDeployer(&_Create2Deployer.CallOpts, salt, codeHash, deployer)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Create2Deployer *Create2DeployerCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Create2Deployer.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Create2Deployer *Create2DeployerSession) Owner() (common.Address, error) {
	return _Create2Deployer.Contract.Owner(&_Create2Deployer.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Create2Deployer *Create2DeployerCallerSession) Owner() (common.Address, error) {
	return _Create2Deployer.Contract.Owner(&_Create2Deployer.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_Create2Deployer *Create2DeployerCaller) Paused(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _Create2Deployer.contract.Call(opts, &out, "paused")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_Create2Deployer *Create2DeployerSession) Paused() (bool, error) {
	return _Create2Deployer.Contract.Paused(&_Create2Deployer.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_Create2Deployer *Create2DeployerCallerSession) Paused() (bool, error) {
	return _Create2Deployer.Contract.Paused(&_Create2Deployer.CallOpts)
}

// Deploy is a paid mutator transaction binding the contract method 0x66cfa057.
//
// Solidity: function deploy(uint256 value, bytes32 salt, bytes code) returns()
func (_Create2Deployer *Create2DeployerTransactor) Deploy(opts *bind.TransactOpts, value *big.Int, salt [32]byte, code []byte) (*types.Transaction, error) {
	return _Create2Deployer.contract.Transact(opts, "deploy", value, salt, code)
}

// Deploy is a paid mutator transaction binding the contract method 0x66cfa057.
//
// Solidity: function deploy(uint256 value, bytes32 salt, bytes code) returns()
func (_Create2Deployer *Create2DeployerSession) Deploy(value *big.Int, salt [32]byte, code []byte) (*types.Transaction, error) {
	return _Create2Deployer.Contract.Deploy(&_Create2Deployer.TransactOpts, value, salt, code)
}

// Deploy is a paid mutator transaction binding the contract method 0x66cfa057.
//
// Solidity: function deploy(uint256 value, bytes32 salt, bytes code) returns()
func (_Create2Deployer *Create2DeployerTransactorSession) Deploy(value *big.Int, salt [32]byte, code []byte) (*types.Transaction, error) {
	return _Create2Deployer.Contract.Deploy(&_Create2Deployer.TransactOpts, value, salt, code)
}

// DeployERC1820Implementer is a paid mutator transaction binding the contract method 0x076c37b2.
//
// Solidity: function deployERC1820Implementer(uint256 value, bytes32 salt) returns()
func (_Create2Deployer *Create2DeployerTransactor) DeployERC1820Implementer(opts *bind.TransactOpts, value *big.Int, salt [32]byte) (*types.Transaction, error) {
	return _Create2Deployer.contract.Transact(opts, "deployERC1820Implementer", value, salt)
}

// DeployERC1820Implementer is a paid mutator transaction binding the contract method 0x076c37b2.
//
// Solidity: function deployERC1820Implementer(uint256 value, bytes32 salt) returns()
func (_Create2Deployer *Create2DeployerSession) DeployERC1820Implementer(value *big.Int, salt [32]byte) (*types.Transaction, error) {
	return _Create2Deployer.Contract.DeployERC1820Implementer(&_Create2Deployer.TransactOpts, value, salt)
}

// DeployERC1820Implementer is a paid mutator transaction binding the contract method 0x076c37b2.
//
// Solidity: function deployERC1820Implementer(uint256 value, bytes32 salt) returns()
func (_Create2Deployer *Create2DeployerTransactorSession) DeployERC1820Implementer(value *big.Int, salt [32]byte) (*types.Transaction, error) {
	return _Create2Deployer.Contract.DeployERC1820Implementer(&_Create2Deployer.TransactOpts, value, salt)
}

// KillCreate2Deployer is a paid mutator transaction binding the contract method 0x64470454.
//
// Solidity: function killCreate2Deployer(address payoutAddress) returns()
func (_Create2Deployer *Create2DeployerTransactor) KillCreate2Deployer(opts *bind.TransactOpts, payoutAddress common.Address) (*types.Transaction, error) {
	return _Create2Deployer.contract.Transact(opts, "killCreate2Deployer", payoutAddress)
}

// KillCreate2Deployer is a paid mutator transaction binding the contract method 0x64470454.
//
// Solidity: function killCreate2Deployer(address payoutAddress) returns()
func (_Create2Deployer *Create2DeployerSession) KillCreate2Deployer(payoutAddress common.Address) (*types.Transaction, error) {
	return _Create2Deployer.Contract.KillCreate2Deployer(&_Create2Deployer.TransactOpts, payoutAddress)
}

// KillCreate2Deployer is a paid mutator transaction binding the contract method 0x64470454.
//
// Solidity: function killCreate2Deployer(address payoutAddress) returns()
func (_Create2Deployer *Create2DeployerTransactorSession) KillCreate2Deployer(payoutAddress common.Address) (*types.Transaction, error) {
	return _Create2Deployer.Contract.KillCreate2Deployer(&_Create2Deployer.TransactOpts, payoutAddress)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_Create2Deployer *Create2DeployerTransactor) Pause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Create2Deployer.contract.Transact(opts, "pause")
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_Create2Deployer *Create2DeployerSession) Pause() (*types.Transaction, error) {
	return _Create2Deployer.Contract.Pause(&_Create2Deployer.TransactOpts)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_Create2Deployer *Create2DeployerTransactorSession) Pause() (*types.Transaction, error) {
	return _Create2Deployer.Contract.Pause(&_Create2Deployer.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Create2Deployer *Create2DeployerTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Create2Deployer.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Create2Deployer *Create2DeployerSession) RenounceOwnership() (*types.Transaction, error) {
	return _Create2Deployer.Contract.RenounceOwnership(&_Create2Deployer.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Create2Deployer *Create2DeployerTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _Create2Deployer.Contract.RenounceOwnership(&_Create2Deployer.TransactOpts)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Create2Deployer *Create2DeployerTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _Create2Deployer.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Create2Deployer *Create2DeployerSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Create2Deployer.Contract.TransferOwnership(&_Create2Deployer.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Create2Deployer *Create2DeployerTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Create2Deployer.Contract.TransferOwnership(&_Create2Deployer.TransactOpts, newOwner)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_Create2Deployer *Create2DeployerTransactor) Unpause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Create2Deployer.contract.Transact(opts, "unpause")
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_Create2Deployer *Create2DeployerSession) Unpause() (*types.Transaction, error) {
	return _Create2Deployer.Contract.Unpause(&_Create2Deployer.TransactOpts)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_Create2Deployer *Create2DeployerTransactorSession) Unpause() (*types.Transaction, error) {
	return _Create2Deployer.Contract.Unpause(&_Create2Deployer.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Create2Deployer *Create2DeployerTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Create2Deployer.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Create2Deployer *Create2DeployerSession) Receive() (*types.Transaction, error) {
	return _Create2Deployer.Contract.Receive(&_Create2Deployer.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Create2Deployer *Create2DeployerTransactorSession) Receive() (*types.Transaction, error) {
	return _Create2Deployer.Contract.Receive(&_Create2Deployer.TransactOpts)
}

// Create2DeployerOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the Create2Deployer contract.
type Create2DeployerOwnershipTransferredIterator struct {
	Event *Create2DeployerOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *Create2DeployerOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Create2DeployerOwnershipTransferred)
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
		it.Event = new(Create2DeployerOwnershipTransferred)
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
func (it *Create2DeployerOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Create2DeployerOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Create2DeployerOwnershipTransferred represents a OwnershipTransferred event raised by the Create2Deployer contract.
type Create2DeployerOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Create2Deployer *Create2DeployerFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*Create2DeployerOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Create2Deployer.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &Create2DeployerOwnershipTransferredIterator{contract: _Create2Deployer.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Create2Deployer *Create2DeployerFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *Create2DeployerOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Create2Deployer.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Create2DeployerOwnershipTransferred)
				if err := _Create2Deployer.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_Create2Deployer *Create2DeployerFilterer) ParseOwnershipTransferred(log types.Log) (*Create2DeployerOwnershipTransferred, error) {
	event := new(Create2DeployerOwnershipTransferred)
	if err := _Create2Deployer.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Create2DeployerPausedIterator is returned from FilterPaused and is used to iterate over the raw logs and unpacked data for Paused events raised by the Create2Deployer contract.
type Create2DeployerPausedIterator struct {
	Event *Create2DeployerPaused // Event containing the contract specifics and raw log

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
func (it *Create2DeployerPausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Create2DeployerPaused)
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
		it.Event = new(Create2DeployerPaused)
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
func (it *Create2DeployerPausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Create2DeployerPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Create2DeployerPaused represents a Paused event raised by the Create2Deployer contract.
type Create2DeployerPaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterPaused is a free log retrieval operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_Create2Deployer *Create2DeployerFilterer) FilterPaused(opts *bind.FilterOpts) (*Create2DeployerPausedIterator, error) {

	logs, sub, err := _Create2Deployer.contract.FilterLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return &Create2DeployerPausedIterator{contract: _Create2Deployer.contract, event: "Paused", logs: logs, sub: sub}, nil
}

// WatchPaused is a free log subscription operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_Create2Deployer *Create2DeployerFilterer) WatchPaused(opts *bind.WatchOpts, sink chan<- *Create2DeployerPaused) (event.Subscription, error) {

	logs, sub, err := _Create2Deployer.contract.WatchLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Create2DeployerPaused)
				if err := _Create2Deployer.contract.UnpackLog(event, "Paused", log); err != nil {
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
func (_Create2Deployer *Create2DeployerFilterer) ParsePaused(log types.Log) (*Create2DeployerPaused, error) {
	event := new(Create2DeployerPaused)
	if err := _Create2Deployer.contract.UnpackLog(event, "Paused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Create2DeployerUnpausedIterator is returned from FilterUnpaused and is used to iterate over the raw logs and unpacked data for Unpaused events raised by the Create2Deployer contract.
type Create2DeployerUnpausedIterator struct {
	Event *Create2DeployerUnpaused // Event containing the contract specifics and raw log

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
func (it *Create2DeployerUnpausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Create2DeployerUnpaused)
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
		it.Event = new(Create2DeployerUnpaused)
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
func (it *Create2DeployerUnpausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Create2DeployerUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Create2DeployerUnpaused represents a Unpaused event raised by the Create2Deployer contract.
type Create2DeployerUnpaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterUnpaused is a free log retrieval operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_Create2Deployer *Create2DeployerFilterer) FilterUnpaused(opts *bind.FilterOpts) (*Create2DeployerUnpausedIterator, error) {

	logs, sub, err := _Create2Deployer.contract.FilterLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return &Create2DeployerUnpausedIterator{contract: _Create2Deployer.contract, event: "Unpaused", logs: logs, sub: sub}, nil
}

// WatchUnpaused is a free log subscription operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_Create2Deployer *Create2DeployerFilterer) WatchUnpaused(opts *bind.WatchOpts, sink chan<- *Create2DeployerUnpaused) (event.Subscription, error) {

	logs, sub, err := _Create2Deployer.contract.WatchLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Create2DeployerUnpaused)
				if err := _Create2Deployer.contract.UnpackLog(event, "Unpaused", log); err != nil {
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
func (_Create2Deployer *Create2DeployerFilterer) ParseUnpaused(log types.Log) (*Create2DeployerUnpaused, error) {
	event := new(Create2DeployerUnpaused)
	if err := _Create2Deployer.contract.UnpackLog(event, "Unpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
