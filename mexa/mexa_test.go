package mexa

import (
	"context"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/require"

	"github.com/valist-io/gasless"
	"github.com/valist-io/gasless/test"
)

// Run-time type checking
var _ gasless.Transactor = (*Mexa)(nil)

func TestSendTransaction(t *testing.T) {
	ctx := context.Background()

	eth, err := ethclient.Dial(os.Getenv("RPC_URL"))
	require.NoError(t, err, "Failed to create ethclient")

	mexa, err := NewMexa(ctx, eth, os.Getenv("BICONOMY_API_KEY"))
	require.NoError(t, err, "Failed to create mexa client")

	contract, err := test.NewTest(test.ValistAddress, eth)
	require.NoError(t, err, "Failed to create contract")

	private, err := crypto.GenerateKey()
	require.NoError(t, err, "Failed to generate private key")

	txopts := gasless.TransactOpts(&bind.TransactOpts{})
	tx, err := contract.CreateOrganization(txopts, "test")
	require.NoError(t, err, "Failed to create transaction")

	signer := gasless.NewPrivateKeySigner(private)
	mexa.Transact(ctx, tx, signer, "9659e431-884a-42e6-9b83-c3cb920a76a7")
	require.NoError(t, err, "Failed to send transaction")
}
