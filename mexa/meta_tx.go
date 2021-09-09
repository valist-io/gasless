package mexa

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/ethereum/go-ethereum/common"
)

const SignatureTypeEIP712 = "EIP712_SIGN"

type MetaTxRequest struct {
	To            string        `json:"to"`
	From          string        `json:"from"`
	ApiId         string        `json:"apiId"`
	Params        []interface{} `json:"params"`
	SignatureType string        `json:"signatureType"`
}

type MetaTxResponse struct {
	TxHash common.Hash `json:"txHash"`
	Log    string      `json:"log"`
	Flag   int         `json:"flag"`
}

// MetaTx executes a meta transaction using the given data.
func (m *Mexa) MetaTx(ctx context.Context, data *MetaTxRequest) (*MetaTxResponse, error) {
	body, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, metaTxNativeURL, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("x-api-key", m.apiKey)

	res, err := m.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	replyData, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var reply MetaTxResponse
	if err := json.Unmarshal(replyData, &reply); err != nil {
		return nil, err
	}

	if reply.TxHash == common.HexToHash("0x0") {
		return nil, fmt.Errorf("%s", replyData)
	}

	return &reply, nil
}
