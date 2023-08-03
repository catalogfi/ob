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
		"0x4FDAAe676608f2a768f9c57BFDAeFA7559283316": {
			ID:       wbtc,
			Decimals: 8,
		},
	},
	"ethereum_optimism": {
		"primary": {
			ID:       primary_eth,
			Decimals: 18,
		},
		"0x1C2172d7BC6F299075fB8799081e21DC3e2CF019": {
			ID:       wbtc,
			Decimals: 8,
		},
	},
}
