package gasless

import (
	"crypto/ecdsa"
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/signer/core"
)

// Signer is a function that returns a signature for the given typed data.
type Signer func(core.TypedData) ([]byte, error)

// NewWalletSigner returns a signer that uses an Ethereum wallet provider.
func NewWalletSigner(account accounts.Account, wallet accounts.Wallet) Signer {
	return func(typedData core.TypedData) ([]byte, error) {
		data, err := json.Marshal(typedData)
		if err != nil {
			return nil, err
		}

		return wallet.SignData(account, accounts.MimetypeTypedData, data)
	}
}

// NewPrivateKeySigner returns a signer that uses an ECDSA private key.
func NewPrivateKeySigner(private *ecdsa.PrivateKey) Signer {
	return func(typedData core.TypedData) ([]byte, error) {
		domainSeparator, err := typedData.HashStruct(EIP712Domain, typedData.Domain.Map())
		if err != nil {
			return nil, err
		}

		typedDataHash, err := typedData.HashStruct(typedData.PrimaryType, typedData.Message)
		if err != nil {
			return nil, err
		}

		rawData := []byte(fmt.Sprintf("\x19\x01%s%s", string(domainSeparator), string(typedDataHash)))
		sighash := crypto.Keccak256(rawData)

		signature, err := crypto.Sign(sighash, private)
		if err != nil {
			return nil, err
		}

		// Transform V from 0/1 to 27/28 according to the yellow paper
		signature[64] += 27
		return signature, nil
	}
}
