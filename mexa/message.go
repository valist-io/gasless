package mexa

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/signer/core"
)

type Message struct {
	From          common.Address `json:"from"`
	To            common.Address `json:"to"`
	Token         common.Address `json:"token"`
	TxGas         uint64         `json:"txGas"`
	TokenGasPrice string         `json:"tokenGasPrice"`
	BatchId       *big.Int       `json:"batchId"`
	BatchNonce    *big.Int       `json:"batchNonce"`
	Deadline      *big.Int       `json:"deadline"`
	Data          string         `json:"data"`
}

func (m *Message) TypedData() core.TypedDataMessage {
	return core.TypedDataMessage{
		"from":          m.From.Hex(),
		"to":            m.To.Hex(),
		"token":         m.Token.Hex(),
		"txGas":         hexutil.EncodeUint64(m.TxGas),
		"tokenGasPrice": m.TokenGasPrice,
		"batchId":       m.BatchId.String(),
		"batchNonce":    m.BatchNonce.String(),
		"deadline":      m.Deadline.String(),
		"data":          m.Data,
	}
}
