package coingecko_api

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	"github.com/sirupsen/logrus"
)

type (
	CoinGecko interface {
		GetCurrentPrice(ctx context.Context, params map[string]string) (*CoinGeckoPriceResponse, error)
	}

	coinGecko struct {
		host       string
		httpClient *http.Client
	}
)

func NewCoinGecko(host string) CoinGecko {
	httpClient := &http.Client{}

	return &coinGecko{
		host:       host,
		httpClient: httpClient,
	}
}

func (cg *coinGecko) GetCurrentPrice(ctx context.Context, queryParams map[string]string) (*CoinGeckoPriceResponse, error) {
	const funcName = "[pkg][coingecko_api]GetCurrentPrice"

	u, err := url.Parse(cg.host)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"err": err.Error(),
		}).Errorf("Error Parsing URL: %s", funcName)
		return nil, err
	}

	u.Path = getSimplePricePath

	q := u.Query()
	for key, value := range queryParams {
		q.Set(key, value)
	}
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"err": err.Error(),
		}).Errorf("Error Making Request: %s", funcName)
		return nil, err
	}

	resp, err := cg.httpClient.Do(req)
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

	var coinGeckoPriceResponse CoinGeckoPriceResponse
	if err := json.Unmarshal(responseBody, &coinGeckoPriceResponse); err != nil {
		logrus.WithFields(logrus.Fields{
			"err":       err.Error(),
			"resp_body": string(responseBody),
		}).Errorf("Error Unmarshal: %s", funcName)
		return nil, err
	}

	return &coinGeckoPriceResponse, nil
}
