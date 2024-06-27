# OrderBook

Orderbook is a order matching engine for [garden.finance](https://garden.finance), implemented in Go. Garden Finance is a decentralised exchange which supports atomic swaps, enabling seamless cross-chain bridging. Through this Orderbook API, users can create orders, track the progress of orders, and complete swaps, while market makers can accept orders and complete trades.

Orderbook serves as the intermediary between users and counterparties during swaps. It facilitates transactions by supporting different order types: market orders for immediate execution, limit orders for specific prices(wip), and Dutch auction orders for dynamic price discovery(wip).

## Configuration

Orderbook uses a configuration file(config.json) to manage its settings. Below is the example config used to enable swaps across ethereum_sepolia and bitcoin_testnet.

```
{
  "CONFIG": {
      "Network": {
          "bitcoin_testnet": {
              "RPC": {
                  "mempool": "https://mempool.space/testnet",
                  "blockstream": "https://blockstream.info/testnet/api"
              },
              "Assets": {
                  "primary": {
                      "Oracle": "https://api.coincap.io/v2/assets/bitcoin",
                      "TokenAddress": "primary",
                      "Decimals": 8
                  }
              },
              "Expiry": 144
          },
          "ethereum_sepolia": {
              "EventWindow": 1000,
              "RPC": {
                  "ethrpc": <sepolia rpc>
              },
              "Assets": {
                  "0x130Ff59B75a415d0bcCc2e996acAf27ce70fD5eF": {
                      "Oracle": "https://api.coincap.io/v2/assets/bitcoin",
                      "TokenAddress": "0x3D1e56247033FE191DC789b5838E366dC04b1b73",
                      "StartBlock": 4291000,
                      "Decimals": 8
                  }
              },
              "Expiry": 7200
          },
       },
      "MinTxLimit": "100000",
      "MaxTxLimit": "150000000",
      "DailyLimit": "600000000"
  },
	"PORT":          <port>,
	"PSQL_DB":       <pgsql connection string>,
	"SERVER_SECRET": <server secret for jwt authentication>,
	"FEEHUB_URL":    <feehub url>,
	"PRICE_URL":     "https://quote.garden.finance"

}
```

### Network Configuration :-

Each network configuration under `CONFIG.Network` includes data about a blockchain network that we want to support.

- Supported networks:
  - `ethereum_sepolia`
  - `bitcoin_testnet`
  - `bitcoin`
  - `bitcoin_regtest`
  - `ethereum`
  - `ethereum_localnet`
  - `ethereum_optimism`
  - `ethereum_arbitrum`
  - `ethereum_polygon`
  - `ethereum_avalanche`
  - `ethereum_bnb`
- `RPC`:

  - For Bitcoin networks, include `mempool` and `blockchain` urls.
  - For EVM networks, include `ethrpc`.

- `Assets`:
  - `<asset>`: Atomicswap contract address deployed in this network.
    - `Oracle`: CoinCap URL for price fetching.
    - `TokenAddress`: Token contract address supported by the specified atomicswap contract address.
    - `Decimals`: Token decimals.
- `Expiry`: Atomic swap expiry time in number of blocks.

## Setup

### Prerequisites

- Go 1.19

### Installation (Using Go commands)

1. Clone the repository
   ```shell
   $ git clone https://github.com/catalogfi/orderbook.git
   ```
   ```shell
   $ cd orderbook
   ```
2. Build the project
   ```shell
   $ go build -tags netgo -ldflags '-s -w' -o ./app ./cmd/rest
   ```
3. Run the application
   ```shell
   $ ./app
   ```

### Installation (Using Docker)

```shell
$ make build
```

```shell
$ make run
```

## Contributing

Contributions are welcome! Please open an issue or submit a pull request if you have any improvements or bug fixes.

## Contact

For any questions or support, please contact us at hello@garden.finance.
