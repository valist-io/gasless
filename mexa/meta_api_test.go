package mexa

import (
	"context"
	"math/big"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/require"
)

func TestMetaApiList(t *testing.T) {
	ctx := context.Background()

	eth, err := ethclient.Dial(os.Getenv("RPC_URL"))
	require.NoError(t, err, "Failed to create ethclient")

	mexa, err := NewMexa(ctx, eth, os.Getenv("BICONOMY_API_KEY"), big.NewInt(0))
	require.NoError(t, err, "Failed to create mexa client")

	_, err = mexa.MetaApi(ctx)
	require.NoError(t, err, "Failed to call meta api")
}
