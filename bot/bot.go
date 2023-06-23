package bot

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/susruth/wbtc-garden/blockchain"
	"github.com/susruth/wbtc-garden/model"
	"github.com/susruth/wbtc-garden/rest"
)

type Executor struct {
	client         rest.Client
	store          Store
	privateKeyFile string
	privateKeys    map[string]string
}

func NewExecutor(keyStr string, store Store, client rest.Client) *Executor {
	return &Executor{
		client:         client,
		store:          store,
		privateKeyFile: keyStr,
	}
}

func (executor *Executor) Run() error {
	for {
		data, err := os.ReadFile(executor.privateKeyFile)
		if err != nil {
			fmt.Println(err)
			continue
		}
		if err := json.Unmarshal(data, &executor.privateKeys); err != nil {
			fmt.Println(err)
			continue
		}

		orders, err := executor.client.GetInitiatorInitiateOrders()
		if err != nil {
			fmt.Println(err)
			continue
		}
		for _, order := range orders {
			if err := executor.handleInitiatorInitiateOrder(order); err != nil {
				fmt.Println(err)
				continue
			}
		}

		orders, err = executor.client.GetInitiatorRedeemOrders()
		if err != nil {
			fmt.Println(err)
			continue
		}
		for _, order := range orders {
			secret, err := executor.store.Secret(order.SecretHash)
			if err != nil {
				fmt.Println(err)
				continue
			}

			secretBytes, err := hex.DecodeString(secret)
			if err != nil {
				fmt.Println(err)
				continue
			}
			if err := executor.handleInitiatorRedeemOrder(order, secretBytes); err != nil {
				fmt.Println(err)
				continue
			}
		}

		orders, err = executor.client.GetFollowerInitiateOrders()
		if err != nil {
			fmt.Println(err)
			continue
		}
		for _, order := range orders {
			if err := executor.handleFollowerInitiateOrder(order); err != nil {
				fmt.Println(err)
				continue
			}
		}

		orders, err = executor.client.GetFollowerRedeemOrders()
		if err != nil {
			fmt.Println(err)
			continue
		}
		for _, order := range orders {
			if err := executor.handleFollowerRedeemOrder(order); err != nil {
				fmt.Println(err)
				continue
			}
		}

		time.Sleep(15 * time.Second)
	}
}

func (executor *Executor) handleInitiatorInitiateOrder(order model.Order) error {
	status, err := executor.store.Status(order.SecretHash)
	if err != nil {
		return err
	}
	if status == InitiatorInitiated {
		return nil
	}

	privateKey := executor.privateKeys[string(order.InitiatorAtomicSwap.Chain)]
	if privateKey == "" {
		return fmt.Errorf("private key not found for chain %s", order.InitiatorAtomicSwap.Chain)
	}
	initiatorSwap, err := blockchain.LoadInitiatorSwap(*order.InitiatorAtomicSwap, privateKey, order.SecretHash)
	if err != nil {
		return err
	}
	txHash, err := initiatorSwap.Initiate()
	if err != nil {
		return err
	}
	if err := executor.store.PutStatus(order.SecretHash, InitiatorInitiated); err != nil {
		return err
	}
	fmt.Println("Initiator initiated swap", txHash)
	return nil
}

func (executor *Executor) handleInitiatorRedeemOrder(order model.Order, secret []byte) error {
	status, err := executor.store.Status(order.SecretHash)
	if err != nil {
		return err
	}
	if status == InitiatorRedeemed {
		return nil
	}

	privateKey := executor.privateKeys[string(order.InitiatorAtomicSwap.Chain)]
	if privateKey == "" {
		return fmt.Errorf("private key not found for chain %s", order.InitiatorAtomicSwap.Chain)
	}
	redeemerSwap, err := blockchain.LoadRedeemerSwap(*order.InitiatorAtomicSwap, privateKey, order.SecretHash)
	if err != nil {
		return err
	}
	txHash, err := redeemerSwap.Redeem(secret)
	if err != nil {
		return err
	}

	if err := executor.store.PutStatus(order.SecretHash, InitiatorRedeemed); err != nil {
		return err
	}
	fmt.Println("Initiator redeemed swap", txHash)
	return nil
}

func (executor *Executor) handleFollowerInitiateOrder(order model.Order) error {
	status, err := executor.store.Status(order.SecretHash)
	if err != nil {
		return err
	}
	if status == FollowerInitiated {
		return nil
	}

	privateKey := executor.privateKeys[string(order.FollowerAtomicSwap.Chain)]
	if privateKey == "" {
		return fmt.Errorf("private key not found for chain %s", order.FollowerAtomicSwap.Chain)
	}
	initiatorSwap, err := blockchain.LoadInitiatorSwap(*order.FollowerAtomicSwap, privateKey, order.SecretHash)
	if err != nil {
		return err
	}
	txHash, err := initiatorSwap.Initiate()
	if err != nil {
		return err
	}
	if err := executor.store.PutStatus(order.SecretHash, FollowerInitiated); err != nil {
		return err
	}
	fmt.Println("Follower initiated swap", txHash)
	return nil
}

func (executor *Executor) handleFollowerRedeemOrder(order model.Order) error {
	status, err := executor.store.Status(order.SecretHash)
	if err != nil {
		return err
	}
	if status == FollowerRedeemed {
		return nil
	}

	privateKey := executor.privateKeys[string(order.FollowerAtomicSwap.Chain)]
	if privateKey == "" {
		return fmt.Errorf("private key not found for chain %s", order.FollowerAtomicSwap.Chain)
	}
	redeemerSwap, err := blockchain.LoadRedeemerSwap(*order.FollowerAtomicSwap, privateKey, order.SecretHash)
	if err != nil {
		return err
	}

	secret, err := hex.DecodeString(order.Secret)
	if err != nil {
		return err
	}

	txHash, err := redeemerSwap.Redeem(secret)
	if err != nil {
		return err
	}
	if err := executor.store.PutStatus(order.SecretHash, FollowerRedeemed); err != nil {
		return err
	}
	fmt.Println("Follower redeemed swap", txHash)
	return nil
}
