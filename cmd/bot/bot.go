package bot

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/spf13/cobra"
	"github.com/susruth/wbtc-garden/bot"
	"github.com/susruth/wbtc-garden/store"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// TODO: CREATE ORDER
/**
 * 1.	Setup User Config as Json and take it as Config
 * 2. 	Create a New Order cmd
 * 3.	Create a fill order cmd
 */

// TAKING FOLLWOING PARAMS FROM CONFIG
/**
 * creator ,sendAddress,recieveAddress
 */

type AddressData struct {
	Creator        string `json:"creator"`
	SendAddress    string `json:"sendAddress"`
	ReceiveAddress string `json:"receiveAddress"`
}

func NewCreateCommand(store store.Store) *cobra.Command {
	var (
		configPath    string
		orderPair     string
		sendAmount    string
		recieveAmount string
	)

	var cmd = &cobra.Command{
		Use:   "Create",
		Short: "Create a new order",
		Run: func(c *cobra.Command, args []string) {
			file, err := ioutil.ReadFile(configPath)
			if err != nil {
				cobra.CheckErr(err)
				return
			}

			var addressData AddressData
			err = json.Unmarshal([]byte(file), &addressData)
			if err != nil {
				cobra.CheckErr(fmt.Sprintf("Error while unmarshalling json: %v", err))
				return
			}

			secret := generateRandomString(16)
			hash := sha256.Sum256([]byte(secret))
			secretHash := hex.EncodeToString(hash[:])

			id, err := store.CreateOrder(addressData.Creator, addressData.SendAddress, addressData.ReceiveAddress, orderPair, sendAmount, recieveAmount, secretHash)
			if err != nil {
				cobra.CheckErr(fmt.Sprintf("Error while creating order: %v", err))
				return
			}

			secretStore, err := bot.NewStore(sqlite.Open("secret.db"), &gorm.Config{})
			if err != nil {
				cobra.CheckErr(fmt.Sprintf("Error while creating secret store: %v", err))
				return
			}
			if err = secretStore.PutSecret(secretHash, secret); err != nil {
				cobra.CheckErr(fmt.Sprintf("Error while creating secret store: %v", err))
				return
			}

			fmt.Println("Order created with id: ", id)
		},
	}

	cmd.Flags().StringVar(&configPath, "config-file", "./config.json", "config file (default is ./config.json)")
	cmd.Flags().StringVar(&orderPair, "order-pair", "", "User should provide the order pair")
	cmd.MarkFlagRequired("order-pair")
	cmd.Flags().StringVar(&sendAmount, "send-amount", "", "User should provide the send amount")
	cmd.MarkFlagRequired("send-amount")
	cmd.Flags().StringVar(&recieveAmount, "recieve-amount", "", "User should provide the recieve amount")
	cmd.MarkFlagRequired("recieve-amount")
	return cmd
}

func NewFillCommand(store store.Store) *cobra.Command {
	var (
		configPath string
		orderId    uint
	)
	var cmd = &cobra.Command{
		Use:   "Fill",
		Short: "Fill an order",
		Run: func(c *cobra.Command, args []string) {
			file, err := ioutil.ReadFile(configPath)
			if err != nil {
				cobra.CheckErr(err)
				return
			}

			var addressData AddressData
			err = json.Unmarshal([]byte(file), &addressData)
			if err != nil {
				cobra.CheckErr(fmt.Sprintf("Error while unmarshalling json: %v", err))
				return
			}

			if err := store.FillOrder(orderId, addressData.Creator, addressData.SendAddress, addressData.ReceiveAddress); err != nil {
				cobra.CheckErr(fmt.Sprintf("Error while filling order: %v", err))
				return
			}

			fmt.Println("Order filled successfully")

		}}
	cmd.Flags().StringVar(&configPath, "config-file", "./config.json", "config file (default is ./config.json)")
	cmd.Flags().UintVar(&orderId, "order-id", 0, "User should provide the order id")
	cmd.MarkFlagRequired("order-id")
	return cmd
}

func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	bytes := make([]byte, length)
	rand.Read(bytes)
	for i, b := range bytes {
		bytes[i] = charset[b%byte(len(charset))]
	}
	return string(bytes)
}

// CreateOrder(creator, sendAddress, recieveAddress, orderPair, sendAmount, recieveAmount, secretHash string) (uint, error)
