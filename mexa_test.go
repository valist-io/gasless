package mexa

import (
	"context"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"

	"github.com/valist-io/mexa/forwarder"
)

func TestMexa(t *testing.T) {
	ctx := context.Background()

	private, err := crypto.GenerateKey()
	require.NoError(t, err, "Failed to generate private key")
	address := crypto.PubkeyToAddress(private.PublicKey)

	chainID := big.NewInt(1337)
	backend := backends.NewSimulatedBackend(core.GenesisAlloc{
		address: {Balance: big.NewInt(9223372036854775807)},
	}, 8000029)
	defer backend.Close()

	opts, err := bind.NewKeyedTransactorWithChainID(private, chainID)
	require.NoError(t, err, "Failed to create transactor")

	forwarderAddress, _, forwarderContract, err := forwarder.DeployForwarder(opts, backend, address)
	require.NoError(t, err, "Failed to deploy forwarder contract")

	_, err = forwarderContract.RegisterDomainSeparator(opts, DomainName, DomainVersion)
	require.NoError(t, err, "Failed to register domain separator")

	backend.Commit()

	batch := big.NewInt(0)
	salt := common.LeftPadBytes(chainID.Bytes(), 32)
	deadline := big.NewInt(time.Now().Add(time.Hour).Unix())

	mexa := Mexa{
		salt:     salt,
		address:  forwarderAddress,
		contract: forwarderContract,
	}

	nonce, err := mexa.Nonce(ctx, address, batch)
	require.NoError(t, err, "Failed to get nonce")

	message := forwarder.ERC20ForwardRequestTypesERC20ForwardRequest{
		From:          address,
		To:            common.HexToAddress("0x0"),
		Token:         common.HexToAddress("0x0"),
		TxGas:         big.NewInt(0),
		TokenGasPrice: big.NewInt(0),
		BatchId:       batch,
		BatchNonce:    nonce,
		Deadline:      deadline,
		Data:          []byte("hello"),
	}
	typedData := TypedData(mexa.Domain(), ERC20ForwardRequest(message))

	domainSeparator, err := typedData.HashStruct("EIP712Domain", typedData.Domain.Map())
	require.NoError(t, err, "Failed to get domain separator")

	typedDataHash, err := typedData.HashStruct(typedData.PrimaryType, typedData.Message)
	require.NoError(t, err, "Failed to get typed data hash")

	rawData := []byte(fmt.Sprintf("\x19\x01%s%s", string(domainSeparator), string(typedDataHash)))
	sighash := crypto.Keccak256(rawData)

	signature, err := crypto.Sign(sighash, private)
	require.NoError(t, err, "Failed to sign data")

	// Transform V from 0/1 to 27/28 according to the yellow paper
	signature[64] += 27

	var separator [32]byte
	copy(separator[:], domainSeparator)

	err = forwarderContract.VerifyEIP712(&bind.CallOpts{}, message, separator, signature)
	require.NoError(t, err, "Failed to verify EIP712")
}
