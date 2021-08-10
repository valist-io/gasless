package mexa

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/signer/core"
)

type Message struct {
	ApiId         string
	BatchId       *big.Int
	From          common.Address
	To            common.Address
	Token         common.Address
	TxGas         uint64
	TokenGasPrice *big.Int
	BatchNonce    *big.Int
	Deadline      *big.Int
	Data          []byte
}

type biconomyEIP712Message struct {
	From          string   `json:"from"`
	To            string   `json:"to"`
	Token         string   `json:"token"`
	TxGas         uint64   `json:"txGas"`
	TokenGasPrice string   `json:"tokenGasPrice"`
	BatchId       *big.Int `json:"batchId"`
	BatchNonce    *big.Int `json:"batchNonce"`
	Deadline      string   `json:"deadline"`
	Data          string   `json:"data"`
}

func (m *Message) TypedData() core.TypedDataMessage {
	txGas := hexutil.EncodeUint64(m.TxGas)
	tokenGasPrice := math.HexOrDecimal256(*m.TokenGasPrice)
	batchId := math.HexOrDecimal256(*m.BatchId)
	batchNonce := math.HexOrDecimal256(*m.BatchNonce)
	deadline := math.HexOrDecimal256(*m.Deadline)

	return core.TypedDataMessage{
		"from":          m.From.Hex(),
		"to":            m.To.Hex(),
		"token":         m.Token.Hex(),
		"txGas":         txGas,
		"tokenGasPrice": &tokenGasPrice,
		"batchId":       &batchId,
		"batchNonce":    &batchNonce,
		"deadline":      &deadline,
		"data":          hexutil.Encode(m.Data),
	}
}

func (m *Message) TypedDataJSON() interface{} {
	return biconomyEIP712Message{
		From:          m.From.Hex(),
		To:            m.To.Hex(),
		Token:         m.Token.Hex(),
		TxGas:         m.TxGas,
		TokenGasPrice: m.TokenGasPrice.String(),
		BatchId:       m.BatchId,
		BatchNonce:    m.BatchNonce,
		Deadline:      m.Deadline.String(),
		Data:          hexutil.Encode(m.Data),
	}
}
