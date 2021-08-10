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
	Nonce         *big.Int
	From          common.Address
	To            common.Address
	Token         common.Address
	TxGas         uint64
	TokenGasPrice *big.Int
	BatchNonce    *big.Int
	Deadline      *big.Int
	Data          []byte
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
