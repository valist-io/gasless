package gasless

import (
	"context"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

const EIP712Domain = "EIP712Domain"

// Transactor defines methods for creating meta transactions.
type Transactor interface {
	// Transact creates a meta transaction from the given transaction.
	Transact(context.Context, *types.Transaction, Signer, ...interface{}) (*types.Transaction, error)
}

// NoSign is used to pass a transaction through without signing.
func NoSign(address common.Address, transaction *types.Transaction) (*types.Transaction, error) {
	return transaction, nil
}

// TransactOpts creates new transaction options for a meta transaction from existing options.
func TransactOpts(txopts *bind.TransactOpts) *bind.TransactOpts {
	return &bind.TransactOpts{
		From:      common.HexToAddress("0x0"),
		Nonce:     txopts.Nonce,
		Signer:    NoSign,
		Value:     txopts.Value,
		GasPrice:  txopts.GasPrice,
		GasFeeCap: txopts.GasFeeCap,
		GasTipCap: txopts.GasTipCap,
		GasLimit:  txopts.GasLimit,
		Context:   txopts.Context,
		NoSend:    true,
	}
}
