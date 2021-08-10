package gasless

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/signer/core"
)

const EIP712Domain = "EIP712Domain"

// MetaTransactor is a meta transaction sender.
type MetaTransactor interface {
	// Types returns the typed data types.
	Types() core.Types
	// PrimaryType returns the typed data primary type.
	PrimaryType() string
	// Domain returns the typed data domain.
	Domain() core.TypedDataDomain
	// Nonce returns the latest nonce for the given address.
	Nonce(context.Context, common.Address) (*big.Int, error)
	// SendTransaction submits the meta transaction with the given signature.
	SendTransaction(context.Context, EIP712Message, []byte, []byte) (*types.Transaction, error)
}

// EIP712Message contains meta transaction data.
type EIP712Message interface {
	// TypedData returns the typed data formatted message.
	TypedData() core.TypedDataMessage
}

// SendTransaction signs the given message and submits it to the given Meta provider.
func SendTransaction(ctx context.Context, meta MetaTransactor, message EIP712Message, signer Signer) (*types.Transaction, error) {
	typedData := core.TypedData{
		Types:       meta.Types(),
		PrimaryType: meta.PrimaryType(),
		Domain:      meta.Domain(),
		Message:     message.TypedData(),
	}

	signature, err := signer(typedData)
	if err != nil {
		return nil, err
	}

	domainSeparator, err := typedData.HashStruct(EIP712Domain, typedData.Domain.Map())
	if err != nil {
		return nil, err
	}

	return meta.SendTransaction(ctx, message, domainSeparator, signature)
}
