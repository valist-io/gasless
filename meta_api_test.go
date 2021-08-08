package mexa

import (
	"context"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/require"
)

func TestMetaApiList(t *testing.T) {
	ctx := context.Background()

	eth, err := ethclient.Dial("https://rpc.valist.io")
	require.NoError(t, err, "Failed to create ethclient")

	mexa, err := NewMexa(ctx, eth, os.Getenv("BICONOMY_API_KEY"))
	require.NoError(t, err, "Failed to create mexa client")

	_, err = mexa.MetaApi().List(ctx)
	require.NoError(t, err, "Failed to call meta api")
}
