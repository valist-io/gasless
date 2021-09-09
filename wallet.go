package gasless

import (
	"fmt"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/signer/core"
)

// Wallet is used to sign typed data for meta transactions.
type Wallet struct {
	accounts.Wallet
}

// NewWallet returns a wrapped ethereum wallet.
func NewWallet(wallet accounts.Wallet) Wallet {
	return Wallet{wallet}
}

// SignTypedData returns a signature derived from the typed data.
func (w Wallet) SignTypedData(account accounts.Account, typedData core.TypedData) ([]byte, error) {
	data, err := formatTypedData(typedData)
	if err != nil {
		return nil, err
	}

	signature, err := w.SignData(account, accounts.MimetypeTypedData, data)
	if err != nil {
		return nil, err
	}

	if signature[64] == 0 || signature[64] == 1 {
		signature[64] += 27
	}

	return signature, nil
}

// SignTypedData is identical to SignTypedData, but also takes a passphrase.
func (w Wallet) SignTypedDataWithPassphrase(account accounts.Account, passphrase string, typedData core.TypedData) ([]byte, error) {
	data, err := formatTypedData(typedData)
	if err != nil {
		return nil, err
	}

	signature, err := w.SignDataWithPassphrase(account, passphrase, accounts.MimetypeTypedData, data)
	if err != nil {
		return nil, err
	}

	if signature[64] == 0 || signature[64] == 1 {
		signature[64] += 27
	}

	return signature, nil
}

func formatTypedData(typedData core.TypedData) ([]byte, error) {
	domainSeparator, err := typedData.HashStruct(EIP712Domain, typedData.Domain.Map())
	if err != nil {
		return nil, err
	}

	typedDataHash, err := typedData.HashStruct(typedData.PrimaryType, typedData.Message)
	if err != nil {
		return nil, err
	}

	return []byte(fmt.Sprintf("\x19\x01%s%s", string(domainSeparator), string(typedDataHash))), nil
}
