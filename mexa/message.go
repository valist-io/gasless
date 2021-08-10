package mexa

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	//"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/signer/core"
)

type Message struct {
	ApiId         string         `json:"-"`
	From          common.Address `json:"from"`
	To            common.Address `json:"to"`
	Token         common.Address `json:"token"`
	TxGas         uint64         `json:"txGas"`
	TokenGasPrice *big.Int       `json:"tokenGasPrice"`
	BatchId       *big.Int       `json:"batchId"`
	BatchNonce    *big.Int       `json:"batchNonce"`
	Deadline      *big.Int       `json:"deadline"`
	Data          []byte         `json:"data"`
}

func (m *Message) TypedData() core.TypedDataMessage {
	// txGas := hexutil.EncodeUint64(m.TxGas)
	// tokenGasPrice := math.HexOrDecimal256(*m.TokenGasPrice)
	// batchId := math.HexOrDecimal256(*m.BatchId)
	// batchNonce := math.HexOrDecimal256(*m.BatchNonce)
	// deadline := math.HexOrDecimal256(*m.Deadline)

	return core.TypedDataMessage{
		"from":          m.From.Hex(),
		"to":            m.To.Hex(),
		"token":         m.Token.Hex(),
		"txGas":         hexutil.EncodeUint64(m.TxGas),
		"tokenGasPrice": m.TokenGasPrice.String(),
		"batchId":       m.BatchId.String(),
		"batchNonce":    m.BatchNonce.String(),
		"deadline":      m.Deadline.String(),
		"data":          hexutil.Encode(m.Data),
	}
}
