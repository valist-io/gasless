package mexa

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
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

	abi, err := abi.JSON(strings.NewReader(test.TestABI))
	require.NoError(t, err, "Failed to get contract abi")

	mexa, err := NewMexa(ctx, eth, os.Getenv("BICONOMY_API_KEY"))
	require.NoError(t, err, "Failed to create mexa client")

	private, err := crypto.GenerateKey()
	require.NoError(t, err, "Failed to generate private key")

	builder := gasless.NewMessageBuilder(abi, test.ValistAddress, eth)
	msg, err := builder.Message(ctx, "createOrganization", "test")
	require.NoError(t, err, "Failed to create ethereum message")

	signer := gasless.NewPrivateKeySigner(private)
	mexa.Transact(ctx, msg, signer, "9659e431-884a-42e6-9b83-c3cb920a76a7")
	require.NoError(t, err, "Failed to send transaction")
}
