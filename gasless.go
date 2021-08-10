package gasless

import (
	"context"
	"encoding/json"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/signer/core"
)

type MetaTransactor interface {
	Types() core.Types
	PrimaryType() string
	Domain() core.TypedDataDomain
	Nonce(context.Context, common.Address) (*big.Int, error)
	SendTransaction(context.Context, EIP712Message, []byte, []byte) (*types.Transaction, error)
}

type EIP712Message interface {
	TypedData() core.TypedDataMessage
	TypedDataJSON() interface{} // return TypedData struct that marshals to the correct JSON
}

func SendTransaction(
	ctx context.Context,
	t MetaTransactor,
	message EIP712Message,
	account accounts.Account,
	wallet accounts.Wallet,
) (*types.Transaction, error) {
	typedData := core.TypedData{
		Types:       t.Types(),
		PrimaryType: t.PrimaryType(),
		Domain:      t.Domain(),
		Message:     message.TypedData(),
	}

	data, err := json.Marshal(typedData)
	if err != nil {
		return nil, err
	}

	signature, err := wallet.SignData(account, accounts.MimetypeTypedData, data)
	if err != nil {
		return nil, err
	}

	domainSeparator, err := typedData.HashStruct("EIP712Domain", typedData.Domain.Map())
	if err != nil {
		return nil, err
	}

	return t.SendTransaction(ctx, message, domainSeparator, signature)
}
