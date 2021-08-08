package mexa

import (
	"bytes"
	"context"
	"encoding/json"
	"math/big"
	"net/http"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/signer/core"
)

type MetaTx Mexa

const (
	EIP712SignatureType = "EIP712_SIGN"
	MetaTxDomainName    = "Biconomy Forwarder"
	MetaTxDomainVersion = "1"
)

type MetaTxNativeRequest struct {
	To            string        `json:"to"`
	From          string        `json:"from"`
	ApiId         string        `json:"apiId"`
	Params        []interface{} `json:"params"`
	SignatureType string        `json:"signatureType"`
}

type MetaTxNativeResponse struct {
	TxHash string `json:"txHash"`
	Log    string `json:"log"`
	Flag   int    `json:"flag"`
}

// Native executes a meta transaction using the given data.
func (m *MetaTx) Native(ctx context.Context, data *MetaTxNativeRequest) (*MetaTxNativeResponse, error) {
	body, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, metaTxNativeURL, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("x-api-key", m.key)

	res, err := m.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var reply MetaTxNativeResponse
	if err := json.NewDecoder(res.Body).Decode(&reply); err != nil {
		return nil, err
	}

	return &reply, nil
}

// Domain returns the meta tx typed data domain.
func (m *MetaTx) Domain() core.TypedDataDomain {
	return core.TypedDataDomain{
		Name:              MetaTxDomainName,
		Version:           MetaTxDomainVersion,
		VerifyingContract: m.address.Hex(),
		Salt:              hexutil.Encode(m.salt),
	}
}

// Nonce returns the latest nonce for the given address and batch.
func (m *MetaTx) Nonce(ctx context.Context, address common.Address, batch *big.Int) (*big.Int, error) {
	callopts := bind.CallOpts{
		Context: ctx,
		From:    address,
	}

	return m.contract.GetNonce(&callopts, address, batch)
}
