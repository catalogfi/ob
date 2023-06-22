package bot

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"

	"github.com/susruth/wbtc-garden/blockchain"
	"github.com/susruth/wbtc-garden/model"
	"github.com/susruth/wbtc-garden/rest"
)

type Bot struct {
	client         rest.Client
	privateKeyFile string
	secretFile     string

	secrets     map[string]string
	privateKeys map[string]string
}

func NewBot(secretStr, keyStr string, client rest.Client) *Bot {
	return &Bot{
		secretFile:     secretStr,
		client:         client,
		privateKeyFile: keyStr,
	}
}

func (bot *Bot) Run() error {
	for {
		data, err := os.ReadFile(bot.secretFile)
		if err != nil {
			fmt.Println(err)
			continue
		}
		if err := json.Unmarshal(data, &bot.secrets); err != nil {
			fmt.Println(err)
			continue
		}

		data, err = os.ReadFile(bot.privateKeyFile)
		if err != nil {
			fmt.Println(err)
			continue
		}
		if err := json.Unmarshal(data, &bot.privateKeys); err != nil {
			fmt.Println(err)
			continue
		}

		orders, err := bot.client.GetInitiatorInitiateOrders()
		if err != nil {
			fmt.Println(err)
			continue
		}
		for _, order := range orders {
			if err := bot.handleInitiatorInitiateOrder(order); err != nil {
				fmt.Println(err)
				continue
			}
		}

		orders, err = bot.client.GetInitiatorRedeemOrders()
		if err != nil {
			fmt.Println(err)
			continue
		}
		for _, order := range orders {
			secretBytes, err := hex.DecodeString(bot.secrets[order.SecretHash])
			if err != nil {
				fmt.Println(err)
				continue
			}
			if err := bot.handleInitiatorRedeemOrder(order, secretBytes); err != nil {
				fmt.Println(err)
				continue
			}
		}

		orders, err = bot.client.GetFollowerInitiateOrders()
		if err != nil {
			fmt.Println(err)
			continue
		}
		for _, order := range orders {
			if err := bot.handleFollowerInitiateOrder(order); err != nil {
				fmt.Println(err)
				continue
			}
		}

		orders, err = bot.client.GetFollowerRedeemOrders()
		if err != nil {
			fmt.Println(err)
			continue
		}
		for _, order := range orders {
			if err := bot.handleFollowerRedeemOrder(order); err != nil {
				fmt.Println(err)
				continue
			}
		}

	}
}

func (bot *Bot) handleInitiatorInitiateOrder(order model.Order) error {
	privateKey := bot.privateKeys[string(order.InitiatorAtomicSwap.Chain)]
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
	fmt.Println("Initiator initiated swap", txHash)
	return nil
}

func (bot *Bot) handleInitiatorRedeemOrder(order model.Order, secret []byte) error {
	privateKey := bot.privateKeys[string(order.InitiatorAtomicSwap.Chain)]
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
	fmt.Println("Initiator redeemed swap", txHash)
	return nil
}

func (bot *Bot) handleFollowerInitiateOrder(order model.Order) error {
	privateKey := bot.privateKeys[string(order.FollowerAtomicSwap.Chain)]
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
	fmt.Println("Follower initiated swap", txHash)
	return nil
}

func (bot *Bot) handleFollowerRedeemOrder(order model.Order) error {
	privateKey := bot.privateKeys[string(order.FollowerAtomicSwap.Chain)]
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
	fmt.Println("Follower redeemed swap", txHash)
	return nil
}
