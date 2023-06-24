package cobi

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/susruth/wbtc-garden/model"
	"github.com/susruth/wbtc-garden/rest"
)

func Create(entropy []byte, store Store) *cobra.Command {
	var (
		account       uint32
		url           string
		orderPair     string
		sendAmount    string
		recieveAmount string
	)

	var cmd = &cobra.Command{
		Use:   "create",
		Short: "create a new order",
		Run: func(c *cobra.Command, args []string) {
			secret := [32]byte{}
			if _, err := rand.Read(secret[:]); err != nil {
				cobra.CheckErr(fmt.Sprintf("Error while generating secret: %v", err))
				return
			}
			hash := sha256.Sum256(secret[:])
			secretHash := hex.EncodeToString(hash[:])

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

			fromChain, toChain, _, _, err := model.ParseOrderPair(orderPair)
			if err != nil {
				cobra.CheckErr(fmt.Sprintf("Error while parsing order pair: %v", err))
				return
			}

			fromAddress, err := getAddressString(entropy, fromChain, account, 0)
			if err != nil {
				cobra.CheckErr(fmt.Sprintf("Error while getting address string: %v", err))
				return
			}

			toAddress, err := getAddressString(entropy, toChain, account, 0)
			if err != nil {
				cobra.CheckErr(fmt.Sprintf("Error while getting address string: %v", err))
				return
			}

			id, err := client.CreateOrder(fromAddress, toAddress, orderPair, sendAmount, recieveAmount, secretHash)
			if err != nil {
				cobra.CheckErr(fmt.Sprintf("Error while creating order: %v", err))
				return
			}

			if err = store.PutSecret(secretHash, hex.EncodeToString(secret[:])); err != nil {
				cobra.CheckErr(fmt.Sprintf("Error while creating secret store: %v", err))
				return
			}

			fmt.Println("Order created with id: ", id)
		},
	}

	cmd.Flags().StringVar(&url, "url", "", "URL of the orderbook server")
	cmd.MarkFlagRequired("url")
	cmd.Flags().Uint32Var(&account, "account", 0, "Account to be used (default: 0)")
	cmd.Flags().StringVar(&orderPair, "order-pair", "", "User should provide the order pair")
	cmd.MarkFlagRequired("order-pair")
	cmd.Flags().StringVar(&sendAmount, "send-amount", "", "User should provide the send amount")
	cmd.MarkFlagRequired("send-amount")
	cmd.Flags().StringVar(&recieveAmount, "recieve-amount", "", "User should provide the recieve amount")
	cmd.MarkFlagRequired("recieve-amount")
	return cmd
}
