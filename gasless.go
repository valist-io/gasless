package gasless

import (
	"context"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
)

const EIP712Domain = "EIP712Domain"

// Transactor defines methods for creating meta transactions.
type Transactor interface {
	// Transact creates a meta transaction from the given transaction.
	Transact(context.Context, *ethereum.CallMsg, Signer, ...interface{}) (*types.Transaction, error)
}
