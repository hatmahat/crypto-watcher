package coingecko_api

import (
	"context"
	"crypto-watcher-backend/pkg/http_request"
	"encoding/json"
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
		}).Errorf("%s: Error Parsing URL", funcName)
		return nil, err
	}

	u.Path += getSimplePricePath

	q := u.Query()
	for key, value := range queryParams {
		q.Set(key, value)
	}
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"err": err.Error(),
		}).Errorf("%s: Error Making Request", funcName)
		return nil, err
	}

	responseBody, err := http_request.DoRequest(cg.httpClient.Do, req, funcName)
	if err != nil {
		logrus.WithError(err).Errorf("%s: Error Do Request", funcName)
		return nil, err
	}

	var coinGeckoPriceResponse CoinGeckoPriceResponse
	if err := json.Unmarshal(responseBody, &coinGeckoPriceResponse); err != nil {
		logrus.WithFields(logrus.Fields{
			"err":       err.Error(),
			"resp_body": string(responseBody),
		}).Errorf("%s: Error Unmarshal", funcName)
		return nil, err
	}

	return &coinGeckoPriceResponse, nil
}
