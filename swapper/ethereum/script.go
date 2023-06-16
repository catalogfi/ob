package ethereum

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	solsha3 "github.com/miguelmota/go-solidity-sha3"
	"github.com/susruth/wbtc-garden-server/swapper/ethereum/typings/AtomicSwap"
	"github.com/susruth/wbtc-garden-server/swapper/ethereum/typings/Create2Deployer"
)

func GetAddress(deployerAddr common.Address, callOps *bind.CallOpts, client *ethclient.Client, redeemer, refunder common.Address, secretHash []byte, expiryBlockNumber *big.Int) (common.Address, error) {
	instance, err := Create2Deployer.NewCreate2Deployer(deployerAddr, client)
	if err != nil {
		return *new(common.Address), err
	}
	var salt [32]byte
	copy(salt[:], secretHash)

	deploymentByteCode := getDeploymentByteCode(redeemer, refunder, salt, expiryBlockNumber)
	hash := solsha3.SoliditySHA3(deploymentByteCode)
	var hash32 [32]byte
	copy(hash32[:], hash)
	addr, err := instance.ComputeAddress(callOps, salt, hash32)
	if err != nil {
		panic(err)
	}
	return addr, nil

}

func Deploy(deployerAddr common.Address, auth *bind.TransactOpts, client *ethclient.Client, redeemer, refunder common.Address, secretHash []byte, expiryBlockNumber *big.Int) (common.Hash, error) {
	instance, err := Create2Deployer.NewCreate2Deployer(deployerAddr, client)
	if err != nil {
		return common.Hash{}, err
	}

	var salt [32]byte
	copy(salt[:], secretHash)

	deploymentByteCode := getDeploymentByteCode(redeemer, refunder, salt, expiryBlockNumber)

	tx, err := instance.Deploy(auth, big.NewInt(0), salt, deploymentByteCode)
	if tx == nil {
		return common.Hash{}, err
	}
	bind.WaitMined(context.Background(), client, tx)

	if err != nil {
		return common.Hash{}, err
	}
	fmt.Printf("Deploy tx sent: %s\n", tx.Hash().Hex())
	return tx.Hash(), nil
}

func getDeploymentByteCode(redeemer, refunder common.Address, salt [32]byte, expiryBlockNumber *big.Int) []byte {

	uint256Ty, _ := abi.NewType("uint256", "uint256", nil)
	bytes32Ty, _ := abi.NewType("bytes32", "bytes32", nil)
	addressTy, _ := abi.NewType("address", "address", nil)
	constructorParamsArguments := abi.Arguments{
		{
			Type: addressTy,
		},
		{
			Type: addressTy,
		},
		{
			Type: bytes32Ty,
		},
		{
			Type: uint256Ty,
		},
	}
	constructorParams, _ := constructorParamsArguments.Pack(
		redeemer,
		refunder,
		salt,
		expiryBlockNumber,
	)

	bytecode, _ := hex.DecodeString(AtomicSwap.AtomicSwapBin)
	deploymentByteCode := append(bytecode, constructorParams...)
	return deploymentByteCode

}
