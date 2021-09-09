//go:generate abigen --abi ./forwarder/ForwarderABI.json --bin ./forwarder/Forwarder.bin --pkg forwarder --out ./forwarder/forwarder.go
package mexa

import (
	"context"
	"fmt"
	"math/big"
	"net/http"
	"strings"
	"time"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/signer/core"

	"github.com/valist-io/gasless"
	"github.com/valist-io/gasless/mexa/forwarder"
)

const (
	metaApiUrl          = "https://api.biconomy.io/api/v1/meta-api"
	metaTxNativeURL     = "https://api.biconomy.io/api/v2/meta-tx/native"
	metaTxSystemInfoURL = "https://api.biconomy.io/api/v2/meta-tx/systemInfo"
)

const (
	PrimaryType = "ERC20ForwardRequest"
	Name        = "Biconomy Forwarder"
	Version     = "1"
)

var Types = core.Types{
	"EIP712Domain": []core.Type{
		{Name: "name", Type: "string"},
		{Name: "version", Type: "string"},
		{Name: "verifyingContract", Type: "address"},
		{Name: "salt", Type: "bytes32"},
	},
	"ERC20ForwardRequest": []core.Type{
		{Name: "from", Type: "address"},
		{Name: "to", Type: "address"},
		{Name: "token", Type: "address"},
		{Name: "txGas", Type: "uint256"},
		{Name: "tokenGasPrice", Type: "uint256"},
		{Name: "batchId", Type: "uint256"},
		{Name: "batchNonce", Type: "uint256"},
		{Name: "deadline", Type: "uint256"},
		{Name: "data", Type: "bytes"},
	},
}

// ForwarderAddressMap is a mapping of chain IDs to Biconomy forwarder contract addresses.
var ForwarderAddressMap = map[string]common.Address{
	"1":      common.HexToAddress("0x84a0856b038eaAd1cC7E297cF34A7e72685A8693"), // Ethereum mainnet
	"3":      common.HexToAddress("0x3D1D6A62c588C1Ee23365AF623bdF306Eb47217A"), // Ropsten testnet
	"4":      common.HexToAddress("0xFD4973FeB2031D4409fB57afEE5dF2051b171104"), // Rinkeby testnet
	"5":      common.HexToAddress("0xE041608922d06a4F26C0d4c27d8bCD01daf1f792"), // Goerli testnet
	"42":     common.HexToAddress("0xF82986F574803dfFd9609BE8b9c7B92f63a1410E"), // Kovan testnet
	"56":     common.HexToAddress("0x86C80a8aa58e0A4fa09A69624c31Ab2a6CAD56b8"), // BSC mainnet
	"97":     common.HexToAddress("0x61456BF1715C1415730076BB79ae118E806E74d2"), // Binance testnet
	"137":    common.HexToAddress("0x86C80a8aa58e0A4fa09A69624c31Ab2a6CAD56b8"), // Matic mainnet
	"100":    common.HexToAddress("0x86C80a8aa58e0A4fa09A69624c31Ab2a6CAD56b8"), // xDai mainnet
	"1287":   common.HexToAddress("0x3AF14449e18f2c3677bFCB5F954Dc68d5fb74a75"), // Moonbeam alpha testnet
	"80001":  common.HexToAddress("0x9399BB24DBB5C4b782C70c2969F58716Ebbd6a3b"), // Matic testnet mumbai
	"421611": common.HexToAddress("0x67454E169d613a8e9BA6b06af2D267696EAaAf41"), // Arbitrum test
}

// Mexa implements EIP712 meta transactions.
type Mexa struct {
	eth    *ethclient.Client
	client *http.Client

	abi     abi.ABI
	address common.Address

	apiKey string
	apiID  map[string]string

	chainID *big.Int
	batchID *big.Int

	forwarderAddress  common.Address
	forwarderContract *forwarder.Forwarder
}

// NewTransactor returns a transactor using the given address and json to build meta transaction requests.
func NewTransactor(eth *ethclient.Client, address common.Address, jsonABI string, apiKey string) (*Mexa, error) {
	parsedABI, err := abi.JSON(strings.NewReader(jsonABI))
	if err != nil {
		return nil, err
	}

	chainID, err := eth.ChainID(context.Background())
	if err != nil {
		return nil, err
	}

	forwarderAddress, ok := ForwarderAddressMap[chainID.String()]
	if !ok {
		return nil, fmt.Errorf("Chain ID not supported: %v", chainID)
	}

	forwarderContract, err := forwarder.NewForwarder(forwarderAddress, eth)
	if err != nil {
		return nil, err
	}

	mexa := &Mexa{
		eth:    eth,
		client: &http.Client{},

		abi:     parsedABI,
		address: address,

		apiKey: apiKey,
		apiID:  make(map[string]string),

		chainID: chainID,
		batchID: big.NewInt(0),

		forwarderAddress:  forwarderAddress,
		forwarderContract: forwarderContract,
	}

	res, err := mexa.MetaApi(context.Background())
	if err != nil {
		return nil, err
	}

	for _, api := range res.List {
		mexa.apiID[api.Method] = api.ID
	}

	return mexa, nil
}

func (m *Mexa) Transact(opts *gasless.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	if opts.Context == nil {
		opts.Context = context.Background()
	}

	apiId, ok := m.apiID[method]
	if !ok {
		return nil, fmt.Errorf("ApiID not found for %s", method)
	}

	data, err := m.abi.Pack(method, params...)
	if err != nil {
		return nil, err
	}

	callmsg := ethereum.CallMsg{
		To:   &m.address,
		From: opts.From,
		Data: data,
	}

	gas, err := m.eth.EstimateGas(opts.Context, callmsg)
	if err != nil {
		return nil, err
	}

	callopts := bind.CallOpts{
		Context: opts.Context,
		From:    opts.From,
	}

	nonce, err := m.forwarderContract.GetNonce(&callopts, opts.From, m.batchID)
	if err != nil {
		return nil, err
	}

	message := &Message{
		From:          opts.From,
		To:            m.address,
		Token:         common.HexToAddress("0x0"),
		TxGas:         gas,
		TokenGasPrice: "0",
		BatchId:       m.batchID,
		BatchNonce:    nonce,
		Deadline:      big.NewInt(time.Now().Add(time.Hour).Unix()),
		Data:          hexutil.Encode(data),
	}

	typedData := core.TypedData{
		Types:       Types,
		PrimaryType: PrimaryType,
		Message:     message.TypedData(),
		Domain: core.TypedDataDomain{
			Name:              Name,
			Version:           Version,
			VerifyingContract: m.forwarderAddress.Hex(),
			Salt:              hexutil.Encode(common.LeftPadBytes(m.chainID.Bytes(), 32)),
		},
	}

	signature, err := opts.MetaSigner(opts.From, typedData)
	if err != nil {
		return nil, err
	}

	domainSeparator, err := typedData.HashStruct(gasless.EIP712Domain, typedData.Domain.Map())
	if err != nil {
		return nil, err
	}

	req := &MetaTxRequest{
		ApiId:         apiId,
		To:            message.To.Hex(),
		From:          message.From.Hex(),
		SignatureType: SignatureTypeEIP712,
		Params: []interface{}{
			message,
			hexutil.Encode(domainSeparator),
			hexutil.Encode(signature),
		},
	}

	res, err := m.MetaTx(opts.Context, req)
	if err != nil {
		return nil, fmt.Errorf("transaction failed: %v", err)
	}

	tx, _, err := m.eth.TransactionByHash(opts.Context, res.TxHash)
	return tx, err
}
