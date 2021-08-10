//go:generate abigen --abi ./ValistABI.json --bin ./Valist.bin --pkg test --out ./valist.go
package test

import (
	"github.com/ethereum/go-ethereum/common"
)

var (
	ValistAddress   = common.HexToAddress("0x11FE3bf493FeBc7B2690e6dc2D233b055F53c6E7")
	RegistryAddress = common.HexToAddress("0x3Ab77276F9c78bC7dE414048aa5E5156D217943B")
)
