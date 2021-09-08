package mexa

import (
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/require"

	"github.com/valist-io/gasless"
	"github.com/valist-io/gasless/test"
)

// Run-time type checking
var _ gasless.Transactor = (*Mexa)(nil)

const (
	veryLightScryptN = 2
	veryLightScryptP = 1
	passphrase       = "secret"
)

func TestSendTransaction(t *testing.T) {
	tmp, err := os.MkdirTemp("", "")
	require.NoError(t, err, "Failed to create tmp dir")
	defer os.RemoveAll(tmp)

	signer := keystore.NewKeyStore(tmp, veryLightScryptN, veryLightScryptP)
	account, err := signer.NewAccount(passphrase)
	require.NoError(t, err, "Failed to create account")

	err = signer.Unlock(account, passphrase)
	require.NoError(t, err, "Failed to unlock account")

	eth, err := ethclient.Dial(os.Getenv("RPC_URL"))
	require.NoError(t, err, "Failed to create ethclient")

	mexa, err := NewTransactor(eth, test.ValistAddress, test.TestABI, os.Getenv("BICONOMY_API_KEY"))
	require.NoError(t, err, "Failed to create mexa client")

	wallet := gasless.NewWallet(signer.Wallets()[0])
	txopts := gasless.NewWalletTransactor(account, wallet)

	_, err = mexa.Transact(txopts, "createOrganization", "0xDEADBEEF")
	require.NoError(t, err, "Failed to send transaction")
}
