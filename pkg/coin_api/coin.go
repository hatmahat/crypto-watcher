package coin_api

import (
	"context"
	"crypto-watcher-backend/internal/constant/http_const"
	"crypto-watcher-backend/pkg/http_request"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
)

type (
	Coin interface {
		GetSpecificRate(ctx context.Context, assetIdBase, assetIdQuote string) (*CoinRateResponse, error)
	}

	coin struct {
		host       string
		apiKey     string
		httpClient *http.Client
	}
)

func NewCoin(host, apiKey string, httpClient *http.Client) Coin {
	return &coin{
		host:       host,
		apiKey:     apiKey,
		httpClient: httpClient,
	}
}

func (c *coin) GetSpecificRate(ctx context.Context, assetIdBase, assetIdQuote string) (*CoinRateResponse, error) {
	const funcName = "[pkg][coin_api]GetSpecificRate"

	req, err := http.NewRequest(http.MethodPost, c.host+fmt.Sprintf(getSpecificRatePath, assetIdBase, assetIdQuote), nil)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"err": err.Error(),
		}).Errorf("%s: Error Making Request", funcName)
		return nil, err
	}

	req.Header.Set(http_const.Accept, http_const.ApplicationJson)
	req.Header.Set(http_const.XCoinApiKey, c.apiKey)

	responseBody, err := http_request.DoRequest(c.httpClient.Do, req, funcName)
	if err != nil {
		logrus.WithError(err).Errorf("%s: Error Do Request", funcName)
		return nil, err
	}

	var coinRateResponse CoinRateResponse
	if err := json.Unmarshal(responseBody, &coinRateResponse); err != nil {
		logrus.WithFields(logrus.Fields{
			"err":       err.Error(),
			"resp_body": string(responseBody),
		}).Errorf("%s: Error Unmarshal", funcName)
		return nil, err
	}

	return &coinRateResponse, nil
}
