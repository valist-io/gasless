package mexa

import (
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/signer/core"

	"github.com/valist-io/mexa/forwarder"
)

type ERC20ForwardRequest forwarder.ERC20ForwardRequestTypesERC20ForwardRequest

// TypedData transforms the message into EIP712 conformant typed data.
func TypedData(domain core.TypedDataDomain, message ERC20ForwardRequest) core.TypedData {
	txGas := math.HexOrDecimal256(*message.TxGas)
	tokenGasPrice := math.HexOrDecimal256(*message.TokenGasPrice)
	batchId := math.HexOrDecimal256(*message.BatchId)
	batchNonce := math.HexOrDecimal256(*message.BatchNonce)
	deadline := math.HexOrDecimal256(*message.Deadline)

	return core.TypedData{
		Types: core.Types{
			"EIP712Domain": []core.Type{
				{Name: "name", Type: "string"},
				{Name: "version", Type: "string"},
				{Name: "verifyingContract", Type: "address"},
				{Name: "salt", Type: "bytes32"},
			},
			"ERC20ForwardRequest": []core.Type{
				{Name: "from", Type: "address"},
				{Name: "to", Type: "address"},
				{Name: "token", Type: "address"},
				{Name: "txGas", Type: "uint256"},
				{Name: "tokenGasPrice", Type: "uint256"},
				{Name: "batchId", Type: "uint256"},
				{Name: "batchNonce", Type: "uint256"},
				{Name: "deadline", Type: "uint256"},
				{Name: "data", Type: "bytes"},
			},
		},
		PrimaryType: "ERC20ForwardRequest",
		Domain:      domain,
		Message: core.TypedDataMessage{
			"from":          message.From.Hex(),
			"to":            message.To.Hex(),
			"token":         message.Token.Hex(),
			"txGas":         &txGas,
			"tokenGasPrice": &tokenGasPrice,
			"batchId":       &batchId,
			"batchNonce":    &batchNonce,
			"deadline":      &deadline,
			"data":          hexutil.Encode(message.Data),
		},
	}
}
