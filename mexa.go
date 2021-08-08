package mexa

import (
	"context"
	"fmt"
	"math/big"
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/valist-io/mexa/forwarder"
)

const (
	metaApiUrl          = "https://api.biconomy.io/api/v1/meta-api"
	metaTxNativeURL     = "https://api.biconomy.io/api/v2/meta-tx/native"
	metaTxSystemInfoURL = "https://api.biconomy.io/api/v2/meta-tx/systemInfo" // ?networkId=
)

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

// Mexa is a Biconomy Mexa client.
type Mexa struct {
	key      string
	salt     []byte
	client   *http.Client
	chainID  *big.Int
	address  common.Address
	contract *forwarder.Forwarder
}

// NewMexa returns a new mexa client.
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

	// calculate salt once
	salt := common.LeftPadBytes(chainID.Bytes(), 32)

	return &Mexa{
		key:      key,
		salt:     salt,
		client:   &http.Client{},
		chainID:  chainID,
		address:  address,
		contract: contract,
	}, nil
}

// MetaApi returns the meta api.
func (m *Mexa) MetaApi() *MetaApi {
	api := MetaApi(*m)
	return &api
}

// MetaTx returns the meta tx api.
func (m *Mexa) MetaTx() *MetaTx {
	api := MetaTx(*m)
	return &api
}
