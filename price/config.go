package price

const (
	usdc    = "usd-coin"
	primary = "ethereum"
	wbtc    = "wrapped-bitcoin"
)

type Token struct {
	ID       string
	Decimals uint64
}

var ConfigMap = map[string]map[string]Token{
	"ethereum": {
		"primary": {
			ID:       primary,
			Decimals: 18,
		},
	},
	"ethereum_sepolia": {
		"primary": {
			ID:       primary,
			Decimals: 18,
		},
		"0x6ECd20D2967eD66b88CDb1d5bcF53c6d0497328a": {
			ID:       wbtc,
			Decimals: 8,
		},
	},
}
