package mexa

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/ethereum/go-ethereum/common"
)

type MetaTxLimit struct {
	Type              int     `json:"type"`
	Value             float32 `json:"value"`
	DurationValue     int     `json:"durationValue"`
	DurationUnit      string  `json:"day"`
	LimitStartTime    int64   `json:"limitStartTime"`
	LimitDurationInMs int64   `json:"limitDurationInMs"`
}

type MetaApiInfo struct {
	ContractAddress   common.Address `json:"contractAddress"`
	ID                string         `json:"id"`
	Name              string         `json:"name"`
	URL               string         `json:"url"`
	Version           int            `json:"version"`
	Method            string         `json:"method"`
	MethodType        string         `json:"methodType"`
	ApiType           string         `json:"apiType"`
	MetaTxLimitStatus int            `json:"metaTxLimitStatus"`
	MetaTxLimit       MetaTxLimit    `json:"metaTxLimit"`
}

type MetaApiResponse struct {
	Log   string        `json:"log"`
	Flag  int           `json:"flag"`
	Total int           `json:"total"`
	List  []MetaApiInfo `json:"listApis"`
}

// MetaApi returns a list of all meta tx enabled functions.
func (m *Mexa) MetaApi(ctx context.Context) (*MetaApiResponse, error) {
	req, err := http.NewRequest(http.MethodGet, metaApiUrl, nil)
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

	var reply MetaApiResponse
	if err := json.Unmarshal(replyData, &reply); err != nil {
		return nil, err
	}

	if reply.Flag != 143 {
		return nil, fmt.Errorf("%s", replyData)
	}

	return &reply, nil
}
