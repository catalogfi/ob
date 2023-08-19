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
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"orderId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"secretHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"initiatedAt\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Initiated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"orderId\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"secrectHash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"secret\",\"type\":\"bytes\"}],\"name\":\"Redeemed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"orderId\",\"type\":\"bytes32\"}],\"name\":\"Refunded\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"atomicSwapOrders\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"redeemer\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"initiator\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"expiry\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"initiatedAt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"isFulfilled\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_redeemer\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_expiry\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"_secretHash\",\"type\":\"bytes32\"}],\"name\":\"initiate\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_orderId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"_secret\",\"type\":\"bytes\"}],\"name\":\"redeem\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_orderId\",\"type\":\"bytes32\"}],\"name\":\"refund\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"token\",\"outputs\":[{\"internalType\":\"contractIERC20\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x60a06040523480156200001157600080fd5b506040516200244938038062002449833981810160405281019062000037919062000117565b62000053672dd095bb30bd9e2e60c01b620000aa60201b60201c565b6200006f67326c480190c52ce060c01b620000aa60201b60201c565b8073ffffffffffffffffffffffffffffffffffffffff1660808173ffffffffffffffffffffffffffffffffffffffff16815250505062000149565b50565b600080fd5b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b6000620000df82620000b2565b9050919050565b620000f181620000d2565b8114620000fd57600080fd5b50565b6000815190506200011181620000e6565b92915050565b60006020828403121562000130576200012f620000ad565b5b6000620001408482850162000100565b91505092915050565b6080516122cf6200017a600039600081816104e501528181610ce30152818161119c01526111e901526122cf6000f3fe608060405234801561001057600080fd5b50600436106100575760003560e01c80633f7b9c381461005c5780637249fbb61461009157806397ffc7ae146100ad578063f7ff7207146100c9578063fc0c546a146100e5575b600080fd5b610076600480360381019061007191906115f2565b610103565b60405161008896959493929190611694565b60405180910390f35b6100ab60048036038101906100a691906115f2565b61018c565b005b6100c760048036038101906100c2919061174d565b61052d565b005b6100e360048036038101906100de9190611819565b610d35565b005b6100ed6111e7565b6040516100fa91906118d8565b60405180910390f35b60006020528060005260406000206000915090508060000160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16908060010160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16908060020154908060030154908060040154908060050160009054906101000a900460ff16905086565b6101a067b5e07a6fda7fd95760c01b61120b565b6101b467420fe2ce19cb560360c01b61120b565b6101c86717b296ee4e5a8b6660c01b61120b565b600080600083815260200190815260200160002090506101f267d649cfe6b31e1c3660c01b61120b565b610206673d3f82424f1fd86960c01b61120b565b61021a67a6bf9de09a7ed7db60c01b61120b565b600073ffffffffffffffffffffffffffffffffffffffff168160000160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16036102ad576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016102a490611950565b60405180910390fd5b6102c167f18d4068ca561a4160c01b61120b565b6102d56786481fe2d400662760c01b61120b565b6102e967d1a6cd8c3d69d5f160c01b61120b565b6102fd67a1e9350d78db438b60c01b61120b565b8060050160009054906101000a900460ff161561034f576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610346906119e2565b60405180910390fd5b61036367671c92e2638dac7d60c01b61120b565b610377676f1c1748ba6eb67960c01b61120b565b61038b67968d6a795202704d60c01b61120b565b61039f67e5cf40fbddc92c6060c01b61120b565b43816002015482600301546103b49190611a31565b106103f4576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016103eb90611ab1565b60405180910390fd5b61040867d6e51b285f8f327360c01b61120b565b61041c678aab2151dcb9273260c01b61120b565b60018160050160006101000a81548160ff02191690831515021790555061044d67d120a80c69aacbe260c01b61120b565b6104616766e6a5113879381160c01b61120b565b817ffe509803c09416b28ff3d8f690c8b0c61462a892c46d5430c8fb20abe472daf060405160405180910390a26104a267fee5f9c6ce261e1260c01b61120b565b6104b66787abec9b9c9f3a2360c01b61120b565b6105298160010160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1682600401547f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1661120e9092919063ffffffff16565b5050565b61054167073ca5a32fb6952260c01b61120b565b83338484610559672eb8948d2701f91b60c01b61120b565b61056d67b146933ae23d37bc60c01b61120b565b610581673f5ab2fedf1d3e0460c01b61120b565b61059567528096bf03cda8dd60c01b61120b565b600073ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff1603610604576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016105fb90611b43565b60405180910390fd5b610618674efaeba6ec265baa60c01b61120b565b61062c67609f20ef22df102760c01b61120b565b61064067aced940ef24b86d460c01b61120b565b6106546734503e0a97eaf7e960c01b61120b565b8373ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff16036106c2576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016106b990611bd5565b60405180910390fd5b6106d667c61d1eddd0043fa460c01b61120b565b6106ea67577974ba2856aa4d60c01b61120b565b6106fe67d0e8e3e9290c1bd460c01b61120b565b61071267ca614e050d444f3560c01b61120b565b60008211610755576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161074c90611c67565b60405180910390fd5b61076967821094dbbea0488260c01b61120b565b61077d679df6528a563b7ab760c01b61120b565b61079167b3922217e112fbc660c01b61120b565b6107a567ee2a8f4931fd9fa260c01b61120b565b600081116107e8576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016107df90611cf9565b60405180910390fd5b6107fc679010d900bc78b85c60c01b61120b565b6108106788b9d45dfd7bde8560c01b61120b565b610824679c934fe1b2958afc60c01b61120b565b610838675a51ca728e95530860c01b61120b565b61084c678782d8ae6fbb73a060c01b61120b565b61086067314bd7420f2f889260c01b61120b565b600060028633604051602001610877929190611d28565b6040516020818303038152906040526040516108939190611dc2565b602060405180830381855afa1580156108b0573d6000803e3d6000fd5b5050506040513d601f19601f820116820180604052508101906108d39190611dee565b90506108e967adf0038e3646dfc360c01b61120b565b6108fd67987580a0f1f591ee60c01b61120b565b60008060008381526020019081526020016000206040518060c00160405290816000820160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020016001820160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020016002820154815260200160038201548152602001600482015481526020016005820160009054906101000a900460ff1615151515815250509050610a1767506926a61330241e60c01b61120b565b610a2b67b98984c4789e2d9160c01b61120b565b610a3f672852710e5820c1d660c01b61120b565b600073ffffffffffffffffffffffffffffffffffffffff16816000015173ffffffffffffffffffffffffffffffffffffffff1614610ab2576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610aa990611e67565b60405180910390fd5b610ac667ea6c5f66df7a8bbd60c01b61120b565b610ada67ed68e950390ff32f60c01b61120b565b610aee67bab85a72d02b1d0f60c01b61120b565b60006040518060c001604052808c73ffffffffffffffffffffffffffffffffffffffff1681526020013373ffffffffffffffffffffffffffffffffffffffff1681526020018b81526020014381526020018a8152602001600015158152509050610b6267ade43036fb58a10a60c01b61120b565b8060008085815260200190815260200160002060008201518160000160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555060208201518160010160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555060408201518160020155606082015181600301556080820151816004015560a08201518160050160006101000a81548160ff021916908315150217905550905050610c5867cfffd857c6e45a3460c01b61120b565b610c6c678b142cb5422ee70560c01b61120b565b87837f3dd1f59c2a4b236fc1e76892b9a4b62de617c6a44a56ed208a3ba79c589823ab83606001518460800151604051610ca7929190611e87565b60405180910390a3610cc3678e4352ac609ac85360c01b61120b565b610cd767b4eca4e72252901060c01b61120b565b610d28333083608001517f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16611294909392919063ffffffff16565b5050505050505050505050565b610d4967146c0ae4cc5e141760c01b61120b565b610d5d67394e55641ebfe30060c01b61120b565b610d7167a0d6a7935f6f031160c01b61120b565b60008060008581526020019081526020016000209050610d9b673ee890e58215725c60c01b61120b565b610daf6753163d3c0c144c0b60c01b61120b565b610dc367a034b186fc481f4a60c01b61120b565b600073ffffffffffffffffffffffffffffffffffffffff168160000160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1603610e56576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610e4d90611950565b60405180910390fd5b610e6a67d94cbc40b8cfa66e60c01b61120b565b610e7e6769d0aeb4eb76d7ea60c01b61120b565b610e926747238988f3eebeea60c01b61120b565b610ea667f71e84f464956fed60c01b61120b565b8060050160009054906101000a900460ff1615610ef8576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610eef906119e2565b60405180910390fd5b610f0c67130d9ea0d85f017560c01b61120b565b610f2067768979520e62085f60c01b61120b565b610f346764cc5b89ac80801060c01b61120b565b600060028484604051610f48929190611ee4565b602060405180830381855afa158015610f65573d6000803e3d6000fd5b5050506040513d601f19601f82011682018060405250810190610f889190611dee565b9050610f9e67fba09d6ca7d8ef5560c01b61120b565b610fb267ba55e9166480f2a560c01b61120b565b610fc667f9d568e05005f5b460c01b61120b565b846002828460010160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16604051602001611000929190611d28565b60405160208183030381529060405260405161101c9190611dc2565b602060405180830381855afa158015611039573d6000803e3d6000fd5b5050506040513d601f19601f8201168201806040525081019061105c9190611dee565b1461109c576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161109390611f49565b60405180910390fd5b6110b0676c04d8cc3b26dacc60c01b61120b565b6110c467d80b6667bd5183c860c01b61120b565b60018260050160006101000a81548160ff0219169083151502179055506110f56723cf144203e69e8c60c01b61120b565b611109674c68b800f37e861660c01b61120b565b807f4c9a044220477b4e94dbb0d07ff6ff4ac30d443bef59098c4541b006954778e286868660405161113d93929190611fb8565b60405180910390a26111596763ebf7958b51da1a60c01b61120b565b61116d672f43b2eba3239c7a60c01b61120b565b6111e08260000160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1683600401547f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1661120e9092919063ffffffff16565b5050505050565b7f000000000000000000000000000000000000000000000000000000000000000081565b50565b61128f8363a9059cbb60e01b848460405160240161122d929190611fea565b604051602081830303815290604052907bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19166020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff838183161783525050505061131d565b505050565b611317846323b872dd60e01b8585856040516024016112b593929190612013565b604051602081830303815290604052907bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19166020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff838183161783525050505061131d565b50505050565b600061137f826040518060400160405280602081526020017f5361666545524332303a206c6f772d6c6576656c2063616c6c206661696c65648152508573ffffffffffffffffffffffffffffffffffffffff166113e59092919063ffffffff16565b90506000815114806113a15750808060200190518101906113a09190612076565b5b6113e0576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016113d790612115565b60405180910390fd5b505050565b60606113f484846000856113fd565b90509392505050565b606082471015611442576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401611439906121a7565b60405180910390fd5b6000808673ffffffffffffffffffffffffffffffffffffffff16858760405161146b9190611dc2565b60006040518083038185875af1925050503d80600081146114a8576040519150601f19603f3d011682016040523d82523d6000602084013e6114ad565b606091505b50915091506114be878383876114ca565b92505050949350505050565b6060831561152c576000835103611524576114e48561153f565b611523576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161151a90612213565b60405180910390fd5b5b829050611537565b6115368383611562565b5b949350505050565b6000808273ffffffffffffffffffffffffffffffffffffffff163b119050919050565b6000825111156115755781518083602001fd5b806040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016115a99190612277565b60405180910390fd5b600080fd5b600080fd5b6000819050919050565b6115cf816115bc565b81146115da57600080fd5b50565b6000813590506115ec816115c6565b92915050565b600060208284031215611608576116076115b2565b5b6000611616848285016115dd565b91505092915050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b600061164a8261161f565b9050919050565b61165a8161163f565b82525050565b6000819050919050565b61167381611660565b82525050565b60008115159050919050565b61168e81611679565b82525050565b600060c0820190506116a96000830189611651565b6116b66020830188611651565b6116c3604083018761166a565b6116d0606083018661166a565b6116dd608083018561166a565b6116ea60a0830184611685565b979650505050505050565b6116fe8161163f565b811461170957600080fd5b50565b60008135905061171b816116f5565b92915050565b61172a81611660565b811461173557600080fd5b50565b60008135905061174781611721565b92915050565b60008060008060808587031215611767576117666115b2565b5b60006117758782880161170c565b945050602061178687828801611738565b935050604061179787828801611738565b92505060606117a8878288016115dd565b91505092959194509250565b600080fd5b600080fd5b600080fd5b60008083601f8401126117d9576117d86117b4565b5b8235905067ffffffffffffffff8111156117f6576117f56117b9565b5b602083019150836001820283011115611812576118116117be565b5b9250929050565b600080600060408486031215611832576118316115b2565b5b6000611840868287016115dd565b935050602084013567ffffffffffffffff811115611861576118606115b7565b5b61186d868287016117c3565b92509250509250925092565b6000819050919050565b600061189e6118996118948461161f565b611879565b61161f565b9050919050565b60006118b082611883565b9050919050565b60006118c2826118a5565b9050919050565b6118d2816118b7565b82525050565b60006020820190506118ed60008301846118c9565b92915050565b600082825260208201905092915050565b7f41746f6d6963537761703a206f72646572206e6f7420696e6974617465640000600082015250565b600061193a601e836118f3565b915061194582611904565b602082019050919050565b600060208201905081810360008301526119698161192d565b9050919050565b7f41746f6d6963537761703a206f7264657220616c72656164792066756c66696c60008201527f6c65640000000000000000000000000000000000000000000000000000000000602082015250565b60006119cc6023836118f3565b91506119d782611970565b604082019050919050565b600060208201905081810360008301526119fb816119bf565b9050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b6000611a3c82611660565b9150611a4783611660565b9250828201905080821115611a5f57611a5e611a02565b5b92915050565b7f41746f6d6963537761703a206f72646572206e6f742065787069726564000000600082015250565b6000611a9b601d836118f3565b9150611aa682611a65565b602082019050919050565b60006020820190508181036000830152611aca81611a8e565b9050919050565b7f41746f6d6963537761703a20696e76616c69642072656465656d65722061646460008201527f7265737300000000000000000000000000000000000000000000000000000000602082015250565b6000611b2d6024836118f3565b9150611b3882611ad1565b604082019050919050565b60006020820190508181036000830152611b5c81611b20565b9050919050565b7f41746f6d6963537761703a2072656465656d657220616e6420696e697469617460008201527f6f722063616e6e6f74206265207468652073616d650000000000000000000000602082015250565b6000611bbf6035836118f3565b9150611bca82611b63565b604082019050919050565b60006020820190508181036000830152611bee81611bb2565b9050919050565b7f41746f6d6963537761703a206578706972792073686f756c642062652067726560008201527f61746572207468616e207a65726f000000000000000000000000000000000000602082015250565b6000611c51602e836118f3565b9150611c5c82611bf5565b604082019050919050565b60006020820190508181036000830152611c8081611c44565b9050919050565b7f41746f6d6963537761703a20616d6f756e742063616e6e6f74206265207a657260008201527f6f00000000000000000000000000000000000000000000000000000000000000602082015250565b6000611ce36021836118f3565b9150611cee82611c87565b604082019050919050565b60006020820190508181036000830152611d1281611cd6565b9050919050565b611d22816115bc565b82525050565b6000604082019050611d3d6000830185611d19565b611d4a6020830184611651565b9392505050565b600081519050919050565b600081905092915050565b60005b83811015611d85578082015181840152602081019050611d6a565b60008484015250505050565b6000611d9c82611d51565b611da68185611d5c565b9350611db6818560208601611d67565b80840191505092915050565b6000611dce8284611d91565b915081905092915050565b600081519050611de8816115c6565b92915050565b600060208284031215611e0457611e036115b2565b5b6000611e1284828501611dd9565b91505092915050565b7f41746f6d6963537761703a206475706c6963617465206f726465720000000000600082015250565b6000611e51601b836118f3565b9150611e5c82611e1b565b602082019050919050565b60006020820190508181036000830152611e8081611e44565b9050919050565b6000604082019050611e9c600083018561166a565b611ea9602083018461166a565b9392505050565b82818337600083830152505050565b6000611ecb8385611d5c565b9350611ed8838584611eb0565b82840190509392505050565b6000611ef1828486611ebf565b91508190509392505050565b7f41746f6d6963537761703a20696e76616c696420736563726574000000000000600082015250565b6000611f33601a836118f3565b9150611f3e82611efd565b602082019050919050565b60006020820190508181036000830152611f6281611f26565b9050919050565b600082825260208201905092915050565b6000601f19601f8301169050919050565b6000611f978385611f69565b9350611fa4838584611eb0565b611fad83611f7a565b840190509392505050565b6000604082019050611fcd6000830186611d19565b8181036020830152611fe0818486611f8b565b9050949350505050565b6000604082019050611fff6000830185611651565b61200c602083018461166a565b9392505050565b60006060820190506120286000830186611651565b6120356020830185611651565b612042604083018461166a565b949350505050565b61205381611679565b811461205e57600080fd5b50565b6000815190506120708161204a565b92915050565b60006020828403121561208c5761208b6115b2565b5b600061209a84828501612061565b91505092915050565b7f5361666545524332303a204552433230206f7065726174696f6e20646964206e60008201527f6f74207375636365656400000000000000000000000000000000000000000000602082015250565b60006120ff602a836118f3565b915061210a826120a3565b604082019050919050565b6000602082019050818103600083015261212e816120f2565b9050919050565b7f416464726573733a20696e73756666696369656e742062616c616e636520666f60008201527f722063616c6c0000000000000000000000000000000000000000000000000000602082015250565b60006121916026836118f3565b915061219c82612135565b604082019050919050565b600060208201905081810360008301526121c081612184565b9050919050565b7f416464726573733a2063616c6c20746f206e6f6e2d636f6e7472616374000000600082015250565b60006121fd601d836118f3565b9150612208826121c7565b602082019050919050565b6000602082019050818103600083015261222c816121f0565b9050919050565b600081519050919050565b600061224982612233565b61225381856118f3565b9350612263818560208601611d67565b61226c81611f7a565b840191505092915050565b60006020820190508181036000830152612291818461223e565b90509291505056fea26469706673582212204d3955dac7dcbad30472eb5040637cd76e2f1ec69432ef292770bd2a9fa7e9b864736f6c63430008120033",
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

// AtomicSwapOrders is a free data retrieval call binding the contract method 0x3f7b9c38.
//
// Solidity: function atomicSwapOrders(bytes32 ) view returns(address redeemer, address initiator, uint256 expiry, uint256 initiatedAt, uint256 amount, bool isFulfilled)
func (_AtomicSwap *AtomicSwapCaller) AtomicSwapOrders(opts *bind.CallOpts, arg0 [32]byte) (struct {
	Redeemer    common.Address
	Initiator   common.Address
	Expiry      *big.Int
	InitiatedAt *big.Int
	Amount      *big.Int
	IsFulfilled bool
}, error) {
	var out []interface{}
	err := _AtomicSwap.contract.Call(opts, &out, "atomicSwapOrders", arg0)

	outstruct := new(struct {
		Redeemer    common.Address
		Initiator   common.Address
		Expiry      *big.Int
		InitiatedAt *big.Int
		Amount      *big.Int
		IsFulfilled bool
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Redeemer = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.Initiator = *abi.ConvertType(out[1], new(common.Address)).(*common.Address)
	outstruct.Expiry = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.InitiatedAt = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.Amount = *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)
	outstruct.IsFulfilled = *abi.ConvertType(out[5], new(bool)).(*bool)

	return *outstruct, err

}

// AtomicSwapOrders is a free data retrieval call binding the contract method 0x3f7b9c38.
//
// Solidity: function atomicSwapOrders(bytes32 ) view returns(address redeemer, address initiator, uint256 expiry, uint256 initiatedAt, uint256 amount, bool isFulfilled)
func (_AtomicSwap *AtomicSwapSession) AtomicSwapOrders(arg0 [32]byte) (struct {
	Redeemer    common.Address
	Initiator   common.Address
	Expiry      *big.Int
	InitiatedAt *big.Int
	Amount      *big.Int
	IsFulfilled bool
}, error) {
	return _AtomicSwap.Contract.AtomicSwapOrders(&_AtomicSwap.CallOpts, arg0)
}

// AtomicSwapOrders is a free data retrieval call binding the contract method 0x3f7b9c38.
//
// Solidity: function atomicSwapOrders(bytes32 ) view returns(address redeemer, address initiator, uint256 expiry, uint256 initiatedAt, uint256 amount, bool isFulfilled)
func (_AtomicSwap *AtomicSwapCallerSession) AtomicSwapOrders(arg0 [32]byte) (struct {
	Redeemer    common.Address
	Initiator   common.Address
	Expiry      *big.Int
	InitiatedAt *big.Int
	Amount      *big.Int
	IsFulfilled bool
}, error) {
	return _AtomicSwap.Contract.AtomicSwapOrders(&_AtomicSwap.CallOpts, arg0)
}

// Token is a free data retrieval call binding the contract method 0xfc0c546a.
//
// Solidity: function token() view returns(address)
func (_AtomicSwap *AtomicSwapCaller) Token(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AtomicSwap.contract.Call(opts, &out, "token")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Token is a free data retrieval call binding the contract method 0xfc0c546a.
//
// Solidity: function token() view returns(address)
func (_AtomicSwap *AtomicSwapSession) Token() (common.Address, error) {
	return _AtomicSwap.Contract.Token(&_AtomicSwap.CallOpts)
}

// Token is a free data retrieval call binding the contract method 0xfc0c546a.
//
// Solidity: function token() view returns(address)
func (_AtomicSwap *AtomicSwapCallerSession) Token() (common.Address, error) {
	return _AtomicSwap.Contract.Token(&_AtomicSwap.CallOpts)
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

// Redeem is a paid mutator transaction binding the contract method 0xf7ff7207.
//
// Solidity: function redeem(bytes32 _orderId, bytes _secret) returns()
func (_AtomicSwap *AtomicSwapTransactor) Redeem(opts *bind.TransactOpts, _orderId [32]byte, _secret []byte) (*types.Transaction, error) {
	return _AtomicSwap.contract.Transact(opts, "redeem", _orderId, _secret)
}

// Redeem is a paid mutator transaction binding the contract method 0xf7ff7207.
//
// Solidity: function redeem(bytes32 _orderId, bytes _secret) returns()
func (_AtomicSwap *AtomicSwapSession) Redeem(_orderId [32]byte, _secret []byte) (*types.Transaction, error) {
	return _AtomicSwap.Contract.Redeem(&_AtomicSwap.TransactOpts, _orderId, _secret)
}

// Redeem is a paid mutator transaction binding the contract method 0xf7ff7207.
//
// Solidity: function redeem(bytes32 _orderId, bytes _secret) returns()
func (_AtomicSwap *AtomicSwapTransactorSession) Redeem(_orderId [32]byte, _secret []byte) (*types.Transaction, error) {
	return _AtomicSwap.Contract.Redeem(&_AtomicSwap.TransactOpts, _orderId, _secret)
}

// Refund is a paid mutator transaction binding the contract method 0x7249fbb6.
//
// Solidity: function refund(bytes32 _orderId) returns()
func (_AtomicSwap *AtomicSwapTransactor) Refund(opts *bind.TransactOpts, _orderId [32]byte) (*types.Transaction, error) {
	return _AtomicSwap.contract.Transact(opts, "refund", _orderId)
}

// Refund is a paid mutator transaction binding the contract method 0x7249fbb6.
//
// Solidity: function refund(bytes32 _orderId) returns()
func (_AtomicSwap *AtomicSwapSession) Refund(_orderId [32]byte) (*types.Transaction, error) {
	return _AtomicSwap.Contract.Refund(&_AtomicSwap.TransactOpts, _orderId)
}

// Refund is a paid mutator transaction binding the contract method 0x7249fbb6.
//
// Solidity: function refund(bytes32 _orderId) returns()
func (_AtomicSwap *AtomicSwapTransactorSession) Refund(_orderId [32]byte) (*types.Transaction, error) {
	return _AtomicSwap.Contract.Refund(&_AtomicSwap.TransactOpts, _orderId)
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
	OrderId     [32]byte
	SecretHash  [32]byte
	InitiatedAt *big.Int
	Amount      *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterInitiated is a free log retrieval operation binding the contract event 0x3dd1f59c2a4b236fc1e76892b9a4b62de617c6a44a56ed208a3ba79c589823ab.
//
// Solidity: event Initiated(bytes32 indexed orderId, bytes32 indexed secretHash, uint256 initiatedAt, uint256 amount)
func (_AtomicSwap *AtomicSwapFilterer) FilterInitiated(opts *bind.FilterOpts, orderId [][32]byte, secretHash [][32]byte) (*AtomicSwapInitiatedIterator, error) {

	var orderIdRule []interface{}
	for _, orderIdItem := range orderId {
		orderIdRule = append(orderIdRule, orderIdItem)
	}
	var secretHashRule []interface{}
	for _, secretHashItem := range secretHash {
		secretHashRule = append(secretHashRule, secretHashItem)
	}

	logs, sub, err := _AtomicSwap.contract.FilterLogs(opts, "Initiated", orderIdRule, secretHashRule)
	if err != nil {
		return nil, err
	}
	return &AtomicSwapInitiatedIterator{contract: _AtomicSwap.contract, event: "Initiated", logs: logs, sub: sub}, nil
}

// WatchInitiated is a free log subscription operation binding the contract event 0x3dd1f59c2a4b236fc1e76892b9a4b62de617c6a44a56ed208a3ba79c589823ab.
//
// Solidity: event Initiated(bytes32 indexed orderId, bytes32 indexed secretHash, uint256 initiatedAt, uint256 amount)
func (_AtomicSwap *AtomicSwapFilterer) WatchInitiated(opts *bind.WatchOpts, sink chan<- *AtomicSwapInitiated, orderId [][32]byte, secretHash [][32]byte) (event.Subscription, error) {

	var orderIdRule []interface{}
	for _, orderIdItem := range orderId {
		orderIdRule = append(orderIdRule, orderIdItem)
	}
	var secretHashRule []interface{}
	for _, secretHashItem := range secretHash {
		secretHashRule = append(secretHashRule, secretHashItem)
	}

	logs, sub, err := _AtomicSwap.contract.WatchLogs(opts, "Initiated", orderIdRule, secretHashRule)
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

// ParseInitiated is a log parse operation binding the contract event 0x3dd1f59c2a4b236fc1e76892b9a4b62de617c6a44a56ed208a3ba79c589823ab.
//
// Solidity: event Initiated(bytes32 indexed orderId, bytes32 indexed secretHash, uint256 initiatedAt, uint256 amount)
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
	OrderId     [32]byte
	SecrectHash [32]byte
	Secret      []byte
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterRedeemed is a free log retrieval operation binding the contract event 0x4c9a044220477b4e94dbb0d07ff6ff4ac30d443bef59098c4541b006954778e2.
//
// Solidity: event Redeemed(bytes32 orderId, bytes32 indexed secrectHash, bytes secret)
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

// WatchRedeemed is a free log subscription operation binding the contract event 0x4c9a044220477b4e94dbb0d07ff6ff4ac30d443bef59098c4541b006954778e2.
//
// Solidity: event Redeemed(bytes32 orderId, bytes32 indexed secrectHash, bytes secret)
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

// ParseRedeemed is a log parse operation binding the contract event 0x4c9a044220477b4e94dbb0d07ff6ff4ac30d443bef59098c4541b006954778e2.
//
// Solidity: event Redeemed(bytes32 orderId, bytes32 indexed secrectHash, bytes secret)
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
	OrderId [32]byte
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRefunded is a free log retrieval operation binding the contract event 0xfe509803c09416b28ff3d8f690c8b0c61462a892c46d5430c8fb20abe472daf0.
//
// Solidity: event Refunded(bytes32 indexed orderId)
func (_AtomicSwap *AtomicSwapFilterer) FilterRefunded(opts *bind.FilterOpts, orderId [][32]byte) (*AtomicSwapRefundedIterator, error) {

	var orderIdRule []interface{}
	for _, orderIdItem := range orderId {
		orderIdRule = append(orderIdRule, orderIdItem)
	}

	logs, sub, err := _AtomicSwap.contract.FilterLogs(opts, "Refunded", orderIdRule)
	if err != nil {
		return nil, err
	}
	return &AtomicSwapRefundedIterator{contract: _AtomicSwap.contract, event: "Refunded", logs: logs, sub: sub}, nil
}

// WatchRefunded is a free log subscription operation binding the contract event 0xfe509803c09416b28ff3d8f690c8b0c61462a892c46d5430c8fb20abe472daf0.
//
// Solidity: event Refunded(bytes32 indexed orderId)
func (_AtomicSwap *AtomicSwapFilterer) WatchRefunded(opts *bind.WatchOpts, sink chan<- *AtomicSwapRefunded, orderId [][32]byte) (event.Subscription, error) {

	var orderIdRule []interface{}
	for _, orderIdItem := range orderId {
		orderIdRule = append(orderIdRule, orderIdItem)
	}

	logs, sub, err := _AtomicSwap.contract.WatchLogs(opts, "Refunded", orderIdRule)
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
// Solidity: event Refunded(bytes32 indexed orderId)
func (_AtomicSwap *AtomicSwapFilterer) ParseRefunded(log types.Log) (*AtomicSwapRefunded, error) {
	event := new(AtomicSwapRefunded)
	if err := _AtomicSwap.contract.UnpackLog(event, "Refunded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
