package cobi

import (
	"crypto/rand"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/susruth/wbtc-garden/model"
	"github.com/tyler-smith/go-bip39"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Run() error {
	var cmd = &cobra.Command{
		Use: "COBI - Catalog Order Book clI",
		Run: func(c *cobra.Command, args []string) {
			c.HelpFunc()(c, args)
		},
		DisableAutoGenTag: true,
	}

	entropy, err := readMnemonic()
	if err != nil {
		return err
	}

	config := model.Config{
		RPC: map[model.Chain]string{
			model.BitcoinRegtest:   "http://localhost:30000",
			model.EthereumLocalnet: "http://localhost:8545",
		},
		DEPLOYERS: map[model.Chain]string{
			model.EthereumLocalnet: "0x13b0D85CcB8bf860b6b79AF3029fCA081AE9beF2",
		},
	}

	store, err := NewStore(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		return err
	}

	cmd.AddCommand(Accounts(entropy))
	cmd.AddCommand(Create(entropy, store))
	cmd.AddCommand(Fill(entropy))
	cmd.AddCommand(Execute(entropy, store, config))
	cmd.AddCommand(Balances(entropy, config))
	// cmd.AddCommand(List())

	if err := cmd.Execute(); err != nil {
		return err
	}
	return nil
}

func readMnemonic() ([]byte, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	if data, err := os.ReadFile(homeDir + "/cob/MNEMONIC"); err == nil {
		return bip39.EntropyFromMnemonic(string(data))
	}
	fmt.Println("Generating new mnemonic")
	entropy := [32]byte{}

	if _, err := rand.Read(entropy[:]); err != nil {
		return nil, err
	}
	mnemonic, err := bip39.NewMnemonic(entropy[:])
	if err != nil {
		return nil, err
	}
	fmt.Println(mnemonic)

	file, err := os.Create(homeDir + "/cob/MNEMONIC")
	if err != nil {
		fmt.Println("error above", err)
		return nil, err
	}
	defer file.Close()

	_, err = file.WriteString(mnemonic)
	if err != nil {
		fmt.Println("error here", err)
		return nil, err
	}
	return entropy[:], nil
}
