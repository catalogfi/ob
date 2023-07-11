package cobi

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/susruth/wbtc-garden/blockchain"
	"github.com/susruth/wbtc-garden/model"
	"github.com/susruth/wbtc-garden/rest"
)

func Execute(entropy []byte, store Store, config model.Config) *cobra.Command {
	var (
		url     string
		account uint32
	)

	var cmd = &cobra.Command{
		Use:   "start",
		Short: "Start the atomic swap executor",
		Run: func(c *cobra.Command, args []string) {
			fmt.Println("check")
			for {
				vals, err := getKeys(entropy, model.Ethereum, account, []uint32{0})
				if err != nil {
					cobra.CheckErr(fmt.Sprintf("Error while getting the signing key: %v", err))
					return
				}
				privKey := vals[0].(*ecdsa.PrivateKey)
				client := rest.NewClient(url, privKey.D.Text(16))
				token, err := client.Login()
				if err != nil {
					cobra.CheckErr(fmt.Sprintf("Error while getting the signing key: %v", err))
					return
				}
				if err := client.SetJwt(token); err != nil {
					cobra.CheckErr(fmt.Sprintf("Error to parse signing key: %v", err))
					return
				}
				fmt.Println("")
				fmt.Println("ORDER")
				fmt.Println("")
				orders, err := client.GetInitiatorInitiateOrders()
				if err != nil {
					fmt.Println(err)
					continue
				}

				for _, order := range orders {
					fmt.Println(order, entropy, account, config, store)
					if err := handleInitiatorInitiateOrder(order, entropy, account, config, store); err != nil {
						fmt.Println(err)
						continue
					}
				}

				orders, err = client.GetFollowerInitiateOrders()
				if err != nil {
					fmt.Println(err)
					continue
				}
				for _, order := range orders {
					if err := handleFollowerInitiateOrder(order, entropy, account, config, store); err != nil {
						fmt.Println(err)
						continue
					}
				}

				orders, err = client.GetInitiatorRedeemOrders()
				if err != nil {
					fmt.Println(err)
					continue
				}
				for _, order := range orders {
					secret, err := store.Secret(order.SecretHash)
					if err != nil {
						fmt.Println(err)
						continue
					}

					secretBytes, err := hex.DecodeString(secret)
					if err != nil {
						fmt.Println(err)
						continue
					}
					if err := handleInitiatorRedeemOrder(order, entropy, account, config, store, secretBytes); err != nil {
						fmt.Println(err)
						continue
					}
				}

				orders, err = client.GetFollowerRedeemOrders()
				if err != nil {
					fmt.Println(err)
					continue
				}
				for _, order := range orders {
					if err := handleFollowerRedeemOrder(order, entropy, account, config, store); err != nil {
						fmt.Println(err)
						continue
					}
				}

				time.Sleep(15 * time.Second)
			}
		},
		DisableAutoGenTag: true,
	}
	cmd.Flags().StringVar(&url, "url", "", "url of the orderbook")
	cmd.MarkFlagRequired("url")
	cmd.Flags().Uint32Var(&account, "account", 0, "account number")
	return cmd
}

func handleInitiatorInitiateOrder(order model.Order, entropy []byte, user uint32, config model.Config, store Store) error {
	fromChain, _, _, _, err := model.ParseOrderPair(order.OrderPair)
	if err != nil {
		return err
	}
	keys, err := getKeys(entropy, fromChain, user, []uint32{0})
	if err != nil {
		return err
	}

	status := store.Status(order.SecretHash)
	if status == InitiatorInitiated {
		return nil
	}
	initiatorSwap, err := blockchain.LoadInitiatorSwap(*order.InitiatorAtomicSwap, keys[0], order.SecretHash, config.RPC)
	if err != nil {
		return err
	}
	txHash, err := initiatorSwap.Initiate()
	if err != nil {
		return err
	}
	if err := store.PutStatus(order.SecretHash, InitiatorInitiated); err != nil {
		return err
	}
	fmt.Println("Initiator initiated swap", txHash)
	return nil
}

func handleInitiatorRedeemOrder(order model.Order, entropy []byte, user uint32, config model.Config, store Store, secret []byte) error {
	_, toChain, _, _, err := model.ParseOrderPair(order.OrderPair)
	if err != nil {
		return err
	}
	keys, err := getKeys(entropy, toChain, user, []uint32{0})
	if err != nil {
		return err
	}

	status := store.Status(order.SecretHash)
	if status == InitiatorRedeemed {
		return nil
	}

	redeemerSwap, err := blockchain.LoadRedeemerSwap(*order.FollowerAtomicSwap, keys[0], order.SecretHash, config.RPC)
	if err != nil {
		return err
	}
	txHash, err := redeemerSwap.Redeem(secret)
	if err != nil {
		return err
	}

	if err := store.PutStatus(order.SecretHash, InitiatorRedeemed); err != nil {
		return err
	}
	fmt.Println("Initiator redeemed swap", txHash)
	return nil
}

func handleFollowerInitiateOrder(order model.Order, entropy []byte, user uint32, config model.Config, store Store) error {
	_, toChain, _, _, err := model.ParseOrderPair(order.OrderPair)
	if err != nil {
		return err
	}
	keys, err := getKeys(entropy, toChain, user, []uint32{0})
	if err != nil {
		return err
	}

	status := store.Status(order.SecretHash)
	if status == FollowerInitiated {
		return nil
	}

	initiatorSwap, err := blockchain.LoadInitiatorSwap(*order.FollowerAtomicSwap, keys[0], order.SecretHash, config.RPC)
	if err != nil {
		return err
	}
	txHash, err := initiatorSwap.Initiate()
	if err != nil {
		return err
	}
	if err := store.PutStatus(order.SecretHash, FollowerInitiated); err != nil {
		return err
	}
	fmt.Println("Follower initiated swap", txHash)
	return nil
}

func handleFollowerRedeemOrder(order model.Order, entropy []byte, user uint32, config model.Config, store Store) error {
	fromChain, _, _, _, err := model.ParseOrderPair(order.OrderPair)
	if err != nil {
		return err
	}
	keys, err := getKeys(entropy, fromChain, user, []uint32{0})
	if err != nil {
		return err
	}

	status := store.Status(order.SecretHash)
	if status == FollowerRedeemed {
		return nil
	}

	redeemerSwap, err := blockchain.LoadRedeemerSwap(*order.InitiatorAtomicSwap, keys[0], order.SecretHash, config.RPC)
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
	if err := store.PutStatus(order.SecretHash, FollowerRedeemed); err != nil {
		return err
	}
	fmt.Println("Follower redeemed swap", txHash)
	return nil
}
