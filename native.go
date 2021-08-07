package mexa

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
)

const nativeURL = "https://api.biconomy.io/api/v2/meta-tx/native"

const EIP712SignatureType = "EIP712_SIGN"

type NativeRequest struct {
	To            string        `json:"to"`
	From          string        `json:"from"`
	ApiId         string        `json:"apiId"`
	Params        []interface{} `json:"params"`
	SignatureType string        `json:"signatureType"`
}

type NativeResponse struct {
	TxHash string `json:"txHash"`
	Log    string `json:"log"`
	Flag   int    `json:"flag"`
}

func (m *Mexa) Native(ctx context.Context, data *NativeRequest) (*NativeResponse, error) {
	body, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, nativeURL, bytes.NewBuffer(body))
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

	var reply NativeResponse
	if err := json.NewDecoder(res.Body).Decode(&reply); err != nil {
		return nil, err
	}

	return &reply, nil
}
