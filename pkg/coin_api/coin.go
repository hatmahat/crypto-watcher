package coin_api

import (
	"context"
	"crypto-watcher-backend/internal/constant/http_const"
	"encoding/json"
	"fmt"
	"io"
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

func NewCoin(host, apiKey string) Coin {
	httpClient := &http.Client{}

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
		}).Errorf("Error Making Request: %s", funcName)
		return nil, err
	}

	req.Header.Set(http_const.Accept, http_const.ApplicationJson)
	req.Header.Set(http_const.XCoinApiKey, c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"err":       err.Error(),
			"resp_code": resp.StatusCode,
			"resp":      resp,
		}).Errorf("Error Calling API: %s", funcName)
		return nil, err
	}

	defer resp.Body.Close()
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"err":  err.Error(),
			"resp": resp,
		}).Errorf("Failed to Read Response: %s", funcName)
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		logrus.WithFields(logrus.Fields{
			"resp_code": resp.StatusCode,
			"resp_body": string(responseBody),
		}).Errorf("Error Calling API: %s", funcName)
		return nil, fmt.Errorf("server response status: %d", resp.StatusCode)
	}

	var coinRateResponse CoinRateResponse
	if err := json.Unmarshal(responseBody, &coinRateResponse); err != nil {
		logrus.WithFields(logrus.Fields{
			"err":       err.Error(),
			"resp_body": string(responseBody),
		}).Errorf("Error Unmarshal: %s", funcName)
		return nil, err
	}

	return &coinRateResponse, nil
}
