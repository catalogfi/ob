package config

const (
	usdc        = "usd-coin"
	primary_btc = "bitcoin"
	primary_eth = "ethereum"
	wbtc        = "wrapped-bitcoin"
)

type Token struct {
	ID       string
	Decimals uint64
}

var ConfigMap = map[string]map[string]Token{
	"bitcoin": {
		"primary": {
			ID:       primary_btc,
			Decimals: 8,
		},
	},
	"bitcoin_testnet": {
		"primary": {
			ID:       primary_btc,
			Decimals: 8,
		},
	},
	"ethereum": {
		"primary": {
			ID:       primary_eth,
			Decimals: 18,
		},
	},
	"ethereum_sepolia": {
		"primary": {
			ID:       primary_eth,
			Decimals: 18,
		},
		"0x6ECd20D2967eD66b88CDb1d5bcF53c6d0497328a": {
			ID:       wbtc,
			Decimals: 8,
		},
	},
}
