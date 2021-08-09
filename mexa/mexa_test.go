package mexa

import (
	"context"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/valist-io/gasless"
)

// Run-time type checking
var _ gasless.MetaTransactor = (*Mexa)(nil)

func TestMetaTxNonce(t *testing.T) {
	ctx := context.Background()

	eth, err := ethclient.Dial("https://rpc.valist.io")
	require.NoError(t, err, "Failed to create ethclient")

	mexa, err := NewMexa(ctx, eth, os.Getenv("BICONOMY_API_KEY"))
	require.NoError(t, err, "Failed to create mexa client")

	_, err = mexa.Nonce(ctx, common.HexToAddress("0x0"))
	require.NoError(t, err, "Failed to call meta api")
}

func TestMetaTxDomain(t *testing.T) {
	ctx := context.Background()

	eth, err := ethclient.Dial("https://rpc.valist.io")
	require.NoError(t, err, "Failed to create ethclient")

	chainID, err := eth.ChainID(ctx)
	require.NoError(t, err, "Failed to get chain id")

	mexa, err := NewMexa(ctx, eth, os.Getenv("BICONOMY_API_KEY"))
	require.NoError(t, err, "Failed to create mexa client")

	address := AddressMap[chainID.String()]
	salt := common.LeftPadBytes(chainID.Bytes(), 32)

	domain := mexa.Domain()
	assert.Equal(t, "Biconomy Forwarder", domain.Name)
	assert.Equal(t, "1", domain.Version)
	assert.Equal(t, address.Hex(), domain.VerifyingContract)
	assert.Equal(t, hexutil.Encode(salt), domain.Salt)
}
