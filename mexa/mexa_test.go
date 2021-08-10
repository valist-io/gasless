package mexa

import (
	"context"
	"math/big"
	"os"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/require"

	"github.com/valist-io/gasless"
	"github.com/valist-io/gasless/test"
)

// Run-time type checking
var _ gasless.MetaTransactor = (*Mexa)(nil)

func TestSendTransaction(t *testing.T) {
	ctx := context.Background()

	eth, err := ethclient.Dial(os.Getenv("RPC_URL"))
	require.NoError(t, err, "Failed to create ethclient")

	mexa, err := NewMexa(ctx, eth, os.Getenv("BICONOMY_API_KEY"), big.NewInt(0))
	require.NoError(t, err, "Failed to create mexa client")

	contract, err := test.NewTest(test.ValistAddress, eth)
	require.NoError(t, err, "Failed to create contract")

	private, err := crypto.GenerateKey()
	require.NoError(t, err, "Failed to generate private key")

	address := crypto.PubkeyToAddress(private.PublicKey)
	txopts := &bind.TransactOpts{
		NoSend: true,
		From:   common.HexToAddress("0x0"),
		Signer: func(address common.Address, transaction *types.Transaction) (*types.Transaction, error) {
			return transaction, nil
		},
	}

	tx, err := contract.CreateOrganization(txopts, "test")
	require.NoError(t, err, "Failed to create transaction")

	nonce, err := mexa.Nonce(ctx, address)
	require.NoError(t, err, "Failed to get address nonce")

	message := &Message{
		ApiId:         "9659e431-884a-42e6-9b83-c3cb920a76a7",
		From:          address,
		To:            *tx.To(),
		Token:         common.HexToAddress("0x0"),
		TxGas:         tx.Gas(),
		TokenGasPrice: "0",
		BatchId:       big.NewInt(0),
		BatchNonce:    nonce,
		Deadline:      big.NewInt(time.Now().Add(time.Hour).Unix()),
		Data:          hexutil.Encode(tx.Data()),
	}

	signer := gasless.NewPrivateKeySigner(private)
	tx, err = gasless.SendTransaction(ctx, mexa, message, signer)
	require.NoError(t, err, "Failed to send transaction")
}
