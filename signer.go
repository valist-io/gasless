package gasless

import (
	"crypto/ecdsa"
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/signer/core"
)

// Signer defines functions for signing meta transactions.
type Signer interface {
	// Address returns the address of the signer.
	Address() common.Address
	// Sign returns a signature derived from the typed data.
	Sign(core.TypedData) ([]byte, error)
}

// WalletSigner uses an ethereum wallet to sign meta transactions.
type WalletSigner struct {
	account accounts.Account
	wallet  accounts.Wallet
}

// NewWalletSigner returns a signer that uses an Ethereum wallet provider.
func NewWalletSigner(account accounts.Account, wallet accounts.Wallet) Signer {
	return &WalletSigner{account, wallet}
}

// Address returns the address of the signer.
func (s *WalletSigner) Address() common.Address {
	return s.account.Address
}

// Sign returns a signature derived from the typed data.
func (s *WalletSigner) Sign(typedData core.TypedData) ([]byte, error) {
	data, err := json.Marshal(typedData)
	if err != nil {
		return nil, err
	}

	return s.wallet.SignData(s.account, accounts.MimetypeTypedData, data)
}

// PrivateKeySigner uses a private key to sign meta transactions.
type PrivateKeySigner struct {
	private *ecdsa.PrivateKey
}

// NewPrivateKeySigner returns a signer that uses an ECDSA private key.
func NewPrivateKeySigner(private *ecdsa.PrivateKey) Signer {
	return &PrivateKeySigner{private}
}

// Address returns the address of the signer.
func (s *PrivateKeySigner) Address() common.Address {
	return crypto.PubkeyToAddress(s.private.PublicKey)
}

// Sign returns a signature derived from the typed data.
func (s *PrivateKeySigner) Sign(typedData core.TypedData) ([]byte, error) {
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

	signature, err := crypto.Sign(sighash, s.private)
	if err != nil {
		return nil, err
	}

	// Transform V from 0/1 to 27/28 according to the yellow paper
	signature[64] += 27
	return signature, nil
}
