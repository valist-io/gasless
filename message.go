package gasless

import (
	"context"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

// MessageBuilder is used to build meta transaction messages.
type MessageBuilder struct {
	abi     abi.ABI
	eth     *ethclient.Client
	address common.Address
}

func NewMessageBuilder(jsonABI string, address common.Address, eth *ethclient.Client) (*MessageBuilder, error) {
	parsedABI, err := abi.JSON(strings.NewReader(jsonABI))
	if err != nil {
		return nil, err
	}

	return &MessageBuilder{parsedABI, eth, address}, nil
}

// Message constructs a new meta transaction message.
func (b *MessageBuilder) Message(ctx context.Context, from common.Address, method string, params ...interface{}) (*ethereum.CallMsg, error) {
	data, err := b.abi.Pack(method, params...)
	if err != nil {
		return nil, err
	}

	msg := ethereum.CallMsg{
		From: from,
		To:   &b.address,
		Data: data,
	}

	gas, err := b.eth.EstimateGas(ctx, msg)
	if err != nil {
		return nil, err
	}

	price, err := b.eth.SuggestGasPrice(ctx)
	if err != nil {
		return nil, err
	}

	msg.Gas = gas
	msg.GasPrice = price
	return &msg, nil
}
