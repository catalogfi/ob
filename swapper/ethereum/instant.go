package ethereum

// import (
// 	"bytes"
// 	"context"
// 	"crypto/ecdsa"
// 	"encoding/hex"
// 	"encoding/json"
// 	"fmt"
// 	"io"
// 	"math/big"
// 	"net/http"
// 	"strconv"
// 	"strings"

// 	"github.com/ethereum/go-ethereum/accounts/abi"
// 	"github.com/ethereum/go-ethereum/accounts/abi/bind"
// 	"github.com/ethereum/go-ethereum/common"
// 	"github.com/ethereum/go-ethereum/common/hexutil"
// 	"github.com/ethereum/go-ethereum/crypto"
// 	"github.com/ethereum/go-ethereum/ethclient"
// 	"github.com/catalogfi/wbtc-garden/swapper/ethereum/typings/ERC20"
// )

// type InstantClientConfig struct {
// 	Url              string
// 	Entrypoint       string
// 	Implementation   string
// 	Factory          string
// 	PaymasterAndData string
// }

// type instantClient struct {
// 	config        InstantClientConfig
// 	indexerClient Client
// }

// func InstantWalletWrapper(config InstantClientConfig, client Client) Client {
// 	return &instantClient{config: config, indexerClient: client}
// }

// func (client *instantClient) GetTransactOpts(privKey *ecdsa.PrivateKey) *bind.TransactOpts {
// 	return client.indexerClient.GetTransactOpts(privKey)
// }

// func (client *instantClient) GetCallOpts() *bind.CallOpts {
// 	return client.indexerClient.GetCallOpts()
// }

// func (client *instantClient) RedeemAtomicSwap(contract common.Address, auth *bind.TransactOpts, token common.Address, secret []byte) (string, error) {
// 	return client.indexerClient.RedeemAtomicSwap(contract, auth, token, secret)
// }

// func (client *instantClient) RefundAtomicSwap(contract common.Address, auth *bind.TransactOpts, token common.Address) (string, error) {
// 	return client.indexerClient.RefundAtomicSwap(contract, auth, token)
// }

// func (client *instantClient) GetPublicAddress(privKey *ecdsa.PrivateKey) common.Address {
// 	return client.indexerClient.GetPublicAddress(privKey)
// }

// func (client *instantClient) GetProvider() *ethclient.Client {
// 	return client.indexerClient.GetProvider()
// }

// func (client *instantClient) TransferERC20(privKey *ecdsa.PrivateKey, amount *big.Int, tokenAddr common.Address, toAddr common.Address) (string, error) {
// 	return client.transferERC20Instant(privKey, amount, tokenAddr, toAddr)
// }

// func (client *instantClient) IsFinal(txHash string) (bool, error) {
// 	// TODO: check whether it is an instant wallet transaction, if it is return true, nil
// 	panic("not implemented")

// 	return client.indexerClient.IsFinal(txHash)
// }

// func (client *instantClient) GetCurrentBlock() (uint64, error) {
// 	return client.indexerClient.GetCurrentBlock()
// }

// func (client *instantClient) GetERC20Balance(tokenAddr common.Address, address common.Address) (*big.Int, error) {
// 	return client.indexerClient.GetERC20Balance(tokenAddr, address)
// }

// func (c *instantClient) transferERC20Instant(privKey *ecdsa.PrivateKey, amount *big.Int, tokenAddr common.Address, toAddr common.Address) (string, error) {
// 	erc20ABI, err := ERC20.ERC20MetaData.GetAbi()
// 	if err != nil {
// 		return "", err
// 	}
// 	transferCallData, err := erc20ABI.Pack("transfer", toAddr, amount)
// 	if err != nil {
// 		return "", err
// 	}
// 	iwABI, err := abi.JSON(strings.NewReader(`[{"inputs": [{"internalType": "address","name": "dest","type": "address"},{"internalType": "uint256","name": "value","type": "uint256"},{"internalType": "bytes","name": "func","type": "bytes"}],"name": "execute","outputs": [],"stateMutability": "nonpayable","type": "function"},{"inputs": [{"internalType": "address[]","name": "dest","type": "address[]"},{"internalType": "uint256[]","name": "values","type": "uint256[]"},{"internalType": "bytes[]","name": "func","type": "bytes[]"}],"name": "executeBatch","outputs": [],"stateMutability": "nonpayable","type": "function"}]`))
// 	if err != nil {
// 		return "", err
// 	}

// 	userAddr := crypto.PubkeyToAddress(privKey.PublicKey)
// 	callData, err := iwABI.Pack("execute", tokenAddr, big.NewInt(0), transferCallData)
// 	if err != nil {
// 		return "", err
// 	}
// 	req, err := c.sign(privKey, callData)
// 	if err != nil {
// 		return "", err
// 	}
// 	if _, err := c.executeRequest("/createExecuteOpSignature", req); err != nil {
// 		return "", err
// 	}

// 	wallet, err := c.getWallet(userAddr)
// 	return wallet.RedeemTxDetails.TxHash, err
// }

// func (c *instantClient) executeRequest(path string, req interface{}) (InstantWallet, error) {
// 	buf := new(bytes.Buffer)
// 	if err := json.NewEncoder(buf).Encode(req); err != nil {
// 		return InstantWallet{}, err
// 	}
// 	resp, err := http.Post(c.config.Url+path, "application/json", buf)
// 	if err != nil {
// 		return InstantWallet{}, err
// 	}
// 	defer resp.Body.Close()

// 	if resp.StatusCode != http.StatusOK {
// 		errObj := struct {
// 			Error string `json:"error"`
// 		}{}

// 		if err := json.NewDecoder(resp.Body).Decode(&errObj); err != nil {
// 			errMsg, err := io.ReadAll(resp.Body)
// 			if err != nil {
// 				return InstantWallet{}, fmt.Errorf("failed to read the error message %v", err)
// 			}
// 			return InstantWallet{}, fmt.Errorf("failed to decode the error %v", string(errMsg))
// 		}
// 		return InstantWallet{}, fmt.Errorf("request failed with status %v", errObj.Error)
// 	}

// 	var wallet InstantWallet
// 	if err := json.NewDecoder(resp.Body).Decode(&wallet); err != nil {
// 		return InstantWallet{}, err
// 	}
// 	return wallet, nil
// }

// type FormattedUserOp struct {
// 	Sender               string  `json:"sender"`
// 	Nonce                *uint64 `json:"nonce"`
// 	InitCode             string  `json:"initCode"`
// 	CallData             string  `json:"callData"`
// 	CallGasLimit         uint64  `json:"callGasLimit"`
// 	VerificationGasLimit uint64  `json:"verificationGasLimit"`
// 	PreVerificationGas   uint64  `json:"preVerificationGas"`
// 	MaxFeePerGas         uint64  `json:"maxFeePerGas"`
// 	MaxPriorityFeePerGas uint64  `json:"maxPriorityFeePerGas"`
// 	PaymasterAndData     string  `json:"paymasterAndData"`
// 	Signature            string  `json:"signature"`
// }

// func (client *instantClient) sign(signer *ecdsa.PrivateKey, callData []byte) (FormattedUserOp, error) {
// 	chainID, err := client.indexerClient.GetProvider().ChainID(context.Background())
// 	if err != nil {
// 		return FormattedUserOp{}, err
// 	}

// 	gasFee, priorityGasFee, err := estimateGas(INSTANT)
// 	if err != nil {
// 		return FormattedUserOp{}, err
// 	}

// 	args := getAbiArgs()
// 	callGasLimit := 200000
// 	verificationGasLimit := 100000
// 	preVerificationGas := 200000
// 	userAddr := crypto.PubkeyToAddress(signer.PublicKey)

// 	entrypoint := common.HexToAddress(client.config.Entrypoint)
// 	implementation := common.HexToAddress(client.config.Implementation)
// 	factory := common.HexToAddress(client.config.Factory)
// 	paymasterAndData, err := hex.DecodeString(client.config.PaymasterAndData)
// 	if err != nil {
// 		return FormattedUserOp{}, err
// 	}
// 	wallet, err := client.getWallet(userAddr)
// 	if err != nil {
// 		return FormattedUserOp{}, err
// 	}
// 	underwriter := common.HexToAddress(wallet.SystemAddress)

// 	initCode := []byte{}
// 	nonce := big.NewInt(wallet.ExpectedNonce + 1)
// 	if nonce.Int64() == 0 {
// 		initCode, err = generateCallData(userAddr, underwriter, implementation, factory, new(big.Int).SetUint64(wallet.TimeLock))
// 		if err != nil {
// 			return FormattedUserOp{}, err
// 		}
// 	}

// 	packed, err := args.Pack(&struct {
// 		Sender               common.Address
// 		Nonce                *big.Int
// 		InitCode             []byte
// 		CallData             []byte
// 		CallGasLimit         *big.Int
// 		VerificationGasLimit *big.Int
// 		PreVerificationGas   *big.Int
// 		MaxFeePerGas         *big.Int
// 		MaxPriorityFeePerGas *big.Int
// 		PaymasterAndData     []byte
// 		Signature            []byte
// 	}{
// 		common.HexToAddress(wallet.WalletAddress),
// 		nonce,
// 		initCode,
// 		callData,
// 		big.NewInt(int64(callGasLimit)),
// 		big.NewInt(int64(verificationGasLimit)),
// 		big.NewInt(int64(preVerificationGas)),
// 		gasFee,
// 		priorityGasFee,
// 		paymasterAndData,
// 		[]byte{},
// 	})
// 	if err != nil {
// 		return FormattedUserOp{}, err
// 	}

// 	hash := crypto.Keccak256Hash(
// 		crypto.Keccak256(packed),
// 		common.LeftPadBytes(entrypoint.Bytes(), 32),
// 		common.LeftPadBytes(chainID.Bytes(), 32),
// 	)
// 	signature, err := crypto.Sign(hash[:], signer)
// 	if err != nil {
// 		panic(err)
// 	}

// 	n := nonce.Uint64()
// 	return FormattedUserOp{
// 		Sender:               wallet.WalletAddress,
// 		Nonce:                &n,
// 		InitCode:             hexutil.Encode(initCode),
// 		CallData:             hexutil.Encode(callData),
// 		CallGasLimit:         uint64(callGasLimit),
// 		VerificationGasLimit: uint64(verificationGasLimit),
// 		PreVerificationGas:   uint64(preVerificationGas),
// 		MaxFeePerGas:         gasFee.Uint64(),
// 		MaxPriorityFeePerGas: priorityGasFee.Uint64(),
// 		PaymasterAndData:     hexutil.Encode(paymasterAndData),
// 		Signature:            hexutil.Encode(signature),
// 	}, nil
// }

// func (i *instantClient) getWallet(userAddress common.Address) (InstantWallet, error) {
// 	req := struct {
// 		WalletAddress string `json:"walletAddress"`
// 		UserAddress   string `json:"userAddress"`
// 	}{
// 		UserAddress: userAddress.String(),
// 	}
// 	iw, err := i.executeRequest("/getWalletInfo", req)
// 	if err != nil {
// 		return InstantWallet{}, err
// 	}
// 	return iw, nil
// }

// func getAbiArgs() abi.Arguments {
// 	UserOpType, err := abi.NewType("tuple", "op", []abi.ArgumentMarshaling{
// 		{Name: "sender", InternalType: "Sender", Type: "address"},
// 		{Name: "nonce", InternalType: "Nonce", Type: "uint256"},
// 		{Name: "initCode", InternalType: "InitCode", Type: "bytes"},
// 		{Name: "callData", InternalType: "CallData", Type: "bytes"},
// 		{Name: "callGasLimit", InternalType: "CallGasLimit", Type: "uint256"},
// 		{Name: "verificationGasLimit", InternalType: "VerificationGasLimit", Type: "uint256"},
// 		{Name: "preVerificationGas", InternalType: "PreVerificationGas", Type: "uint256"},
// 		{Name: "maxFeePerGas", InternalType: "MaxFeePerGas", Type: "uint256"},
// 		{Name: "maxPriorityFeePerGas", InternalType: "MaxPriorityFeePerGas", Type: "uint256"},
// 		{Name: "paymasterAndData", InternalType: "PaymasterAndData", Type: "bytes"},
// 		{Name: "signature", InternalType: "Signature", Type: "bytes"},
// 	})
// 	if err != nil {
// 		panic(err)
// 	}
// 	return abi.Arguments{
// 		{Name: "UserOp", Type: UserOpType},
// 	}
// }

// type InstantWallet struct {
// 	WalletAddress   string `json:"walletAddress"`
// 	SystemAddress   string `json:"systemAddress"`
// 	RedeemTxDetails struct {
// 		TxHash string `json:"txHash"`
// 	}
// 	ExpectedNonce int64  `json:"expectedNonce"`
// 	LastSeenNonce int64  `json:"lastSeenNonce"`
// 	TimeLock      uint64 `json:"timeLock"`
// }

// func generateCallData(ownerAddr, underwiterAddr, implAddr, factory common.Address, timeLock *big.Int) ([]byte, error) {
// 	hexTimelock := strconv.FormatInt(timeLock.Int64(), 16)
// 	calldataHex := "f1a57e8f" + strings.Repeat("0000", 6) + ownerAddr.String()[2:] + strings.Repeat("0000", 6) + underwiterAddr.String()[2:] + strings.Repeat("0", 64-len(hexTimelock)) + hexTimelock
// 	callData, err := hex.DecodeString(calldataHex)
// 	if err != nil {
// 		return nil, err
// 	}
// 	factoryABI, err := abi.JSON(strings.NewReader("{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_implementationAddr\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"_initData\",\"type\":\"bytes\"}],\"name\":\"createInstantWallet\",\"outputs\":[{\"internalType\":\"contractInstantWallet\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}"))
// 	if err != nil {
// 		return nil, err
// 	}
// 	initCode, err := factoryABI.Pack("createInstantWallet", implAddr, callData)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return append(factory.Bytes(), initCode...), nil
// }

// type TX_SPEED int

// const (
// 	INSTANT TX_SPEED = iota
// 	FAST
// 	STANDARD
// )

// func estimateGas(speed TX_SPEED) (*big.Int, *big.Int, error) {
// 	resp, err := http.Get("https://api.ethgasstation.info/api/fee-estimate")
// 	if err != nil {
// 		return nil, nil, err
// 	}

// 	estimateGasResp := EstimateGasResponse{}
// 	if json.NewDecoder(resp.Body).Decode(&estimateGasResp); err != nil {
// 		return nil, nil, err
// 	}

// 	switch speed {
// 	case INSTANT:
// 		return big.NewInt(int64(estimateGasResp.GasPrice.Instant)), big.NewInt(int64(estimateGasResp.PriorityFee.Instant)), nil
// 	case FAST:
// 		return big.NewInt(int64(estimateGasResp.GasPrice.Fast)), big.NewInt(int64(estimateGasResp.PriorityFee.Fast)), nil
// 	case STANDARD:
// 		return big.NewInt(int64(estimateGasResp.GasPrice.Standard)), big.NewInt(int64(estimateGasResp.PriorityFee.Standard)), nil
// 	default:
// 		return nil, nil, fmt.Errorf("invalid speed: %v", speed)
// 	}
// }

// type EstimateGasResponse struct {
// 	BaseFee     int     `json:"baseFee"`
// 	BlockNumber int     `json:"blockNumber"`
// 	BlockTime   float64 `json:"blockTime"`
// 	GasPrice    struct {
// 		Fast     int `json:"fast"`
// 		Instant  int `json:"instant"`
// 		Standard int `json:"standard"`
// 	} `json:"gasPrice"`
// 	NextBaseFee int `json:"nextBaseFee"`
// 	PriorityFee struct {
// 		Fast     int `json:"fast"`
// 		Instant  int `json:"instant"`
// 		Standard int `json:"standard"`
// 	} `json:"priorityFee"`
// }
