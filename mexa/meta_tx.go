package mexa

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
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
	TxHash string `json:"txHash"`
	Log    string `json:"log"`
	Flag   int    `json:"flag"`
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
	req.Header.Set("x-api-key", m.key)

	res, err := m.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var reply MetaTxResponse
	if err := json.NewDecoder(res.Body).Decode(&reply); err != nil {
		return nil, err
	}

	return &reply, nil
}
