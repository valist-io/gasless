package gasless

import (
	"context"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/signer/core"
)

const EIP712Domain = "EIP712Domain"

// MetaSignerFn is a function used to sign typed data.
type MetaSignerFn func(from common.Address, typedData core.TypedData) ([]byte, error)

// TransactOpts contains transaction settings.
type TransactOpts struct {
	bind.TransactOpts
	MetaSigner MetaSignerFn
}

// Transactor defines methods for creating meta transactions.
type Transactor interface {
	// Transact creates and submits a meta transaction.
	Transact(opts *TransactOpts, method string, params ...interface{}) (*types.Transaction, error)
}

// NewWalletTransactor returns transact opts using a gasless wallet.
func NewWalletTransactor(account accounts.Account, wallet Wallet) *TransactOpts {
	return &TransactOpts{
		TransactOpts: bind.TransactOpts{
			Context: context.Background(),
			From:    account.Address,
		},
		MetaSigner: func(from common.Address, typedData core.TypedData) ([]byte, error) {
			return wallet.SignTypedData(account, typedData)
		},
	}
}
