package gasless

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/signer/core"
)

type MetaTransactor interface {
	Types() core.Types
	PrimaryType() string
	Domain() core.TypedDataDomain
	Nonce(context.Context, common.Address) (*big.Int, error)
}
