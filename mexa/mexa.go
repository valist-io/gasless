//go:generate abigen --abi ./forwarder/ForwarderABI.json --bin ./forwarder/Forwarder.bin --pkg forwarder --out ./forwarder/forwarder.go
package mexa

import (
	"context"
	"fmt"
	"math/big"
	"net/http"
	"time"

	ethereum "github.com/ethereum/go-ethereum"
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

const PrimaryType = "ERC20ForwardRequest"

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

// AddressMap is a mapping of chain IDs to Biconomy forwarder contract addresses.
var AddressMap = map[string]common.Address{
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
	eth      *ethclient.Client
	key      string
	client   *http.Client
	contract *forwarder.Forwarder
	domain   core.TypedDataDomain
}

// NewMexa returns a client with the given eth client, key, and batchID.
func NewMexa(ctx context.Context, eth *ethclient.Client, key string) (*Mexa, error) {
	chainID, err := eth.ChainID(ctx)
	if err != nil {
		return nil, err
	}

	address, ok := AddressMap[chainID.String()]
	if !ok {
		return nil, fmt.Errorf("Chain ID not supported: %s", chainID)
	}

	contract, err := forwarder.NewForwarder(address, eth)
	if err != nil {
		return nil, err
	}

	salt := common.LeftPadBytes(chainID.Bytes(), 32)

	return &Mexa{
		eth:      eth,
		key:      key,
		client:   &http.Client{},
		contract: contract,
		domain: core.TypedDataDomain{
			Name:              "Biconomy Forwarder",
			Version:           "1",
			VerifyingContract: address.Hex(),
			Salt:              hexutil.Encode(salt),
		},
	}, nil
}

// Transact creates a meta transaction from the given transaction.
func (m *Mexa) Transact(ctx context.Context, msg *ethereum.CallMsg, signer gasless.Signer, params ...interface{}) (*types.Transaction, error) {
	// TODO make this a param
	batchID := big.NewInt(0)

	if len(params) == 0 {
		return nil, fmt.Errorf("mexa transact requires ApiId string param")
	}

	apiId, ok := params[0].(string)
	if !ok {
		return nil, fmt.Errorf("mexa transact requires ApiId string param")
	}

	callopts := bind.CallOpts{
		Context: ctx,
		From:    signer.Address(),
	}

	nonce, err := m.contract.GetNonce(&callopts, signer.Address(), batchID)
	if err != nil {
		return nil, err
	}

	message := &Message{
		From:          signer.Address(),
		To:            *msg.To,
		Token:         common.HexToAddress("0x0"),
		TxGas:         msg.Gas,
		TokenGasPrice: "0",
		BatchId:       batchID,
		BatchNonce:    nonce,
		Deadline:      big.NewInt(time.Now().Add(time.Hour).Unix()),
		Data:          hexutil.Encode(msg.Data),
	}

	fmt.Println(m.key)

	typedData := core.TypedData{
		Types:       Types,
		PrimaryType: PrimaryType,
		Domain:      m.domain,
		Message:     message.TypedData(),
	}

	signature, err := signer.Sign(typedData)
	if err != nil {
		return nil, err
	}

	domainSeparator, err := typedData.HashStruct(gasless.EIP712Domain, typedData.Domain.Map())
	if err != nil {
		return nil, err
	}

	req := &MetaTxRequest{
		To:            message.To.Hex(),
		From:          message.From.Hex(),
		ApiId:         apiId,
		SignatureType: SignatureTypeEIP712,
		Params: []interface{}{
			message,
			hexutil.Encode(domainSeparator),
			hexutil.Encode(signature),
		},
	}

	res, err := m.MetaTx(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("transaction failed: %v", err)
	}

	tx, _, err := m.eth.TransactionByHash(ctx, res.TxHash)
	return tx, err
}
