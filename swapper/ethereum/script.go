package ethereum

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/susruth/wbtc-garden/swapper/ethereum/typings/Deployer"
)

func GetAddress(client Client, deployerAddr common.Address, redeemer, refunder common.Address, tokenAddr common.Address, secretHash []byte, expiryPeriod *big.Int, amount *big.Int) (common.Address, error) {
	instance, err := Deployer.NewDeployer(deployerAddr, client.GetProvider())
	if err != nil {
		return *new(common.Address), err
	}
	var salt [32]byte
	copy(salt[:], secretHash)
	addr, err := instance.ComputeAddress(client.GetCallOpts(), redeemer, refunder, tokenAddr, salt, expiryPeriod, amount)
	if err != nil {
		panic(err)
	}
	return addr, nil
}

func Deploy(auth *bind.TransactOpts, client *ethclient.Client, deployerAddr, redeemer, refunder, tokenAddr common.Address, secretHash []byte, expiryPeriod *big.Int, amount *big.Int) (common.Hash, error) {
	instance, err := Deployer.NewDeployer(deployerAddr, client)
	if err != nil {
		return common.Hash{}, err
	}

	var salt [32]byte
	copy(salt[:], secretHash)

	tx, err := instance.Deploy(auth, redeemer, refunder, tokenAddr, salt, expiryPeriod, amount)
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
