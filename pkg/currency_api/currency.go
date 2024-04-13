package currency_api

import (
	"context"
	"crypto-watcher-backend/pkg/http_request"
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/sirupsen/logrus"
)

type (
	Currency interface {
		GetCurrentCurrency(ctx context.Context, params map[string]string) (*CurrencyAPIResponse, error)
	}

	currency struct {
		host       string
		apiKey     string
		httpClient *http.Client
	}
)

func NewCurrency(host, apiKey string) Currency {
	httpClient := &http.Client{}

	return &currency{
		host:       host,
		apiKey:     apiKey,
		httpClient: httpClient,
	}
}

func (c *currency) GetCurrentCurrency(ctx context.Context, queryParams map[string]string) (*CurrencyAPIResponse, error) {
	const funcName = "[pkg][currency_api]GetCurrentCurrency"

	u, err := url.Parse(c.host)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"err": err.Error(),
		}).Errorf("%s: Error Parsing URL", funcName)
		return nil, err
	}

	u.Path += getLatest

	queryParams[APIKey] = c.apiKey
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

	responseBody, err := http_request.DoRequest(c.httpClient.Do, req, funcName)
	if err != nil {
		logrus.WithError(err).Errorf("%s: Error Do Request", funcName)
		return nil, err
	}

	var currency CurrencyAPIResponse
	if err := json.Unmarshal(responseBody, &currency); err != nil {
		logrus.WithFields(logrus.Fields{
			"err":       err.Error(),
			"resp_body": string(responseBody),
		}).Errorf("%s: Error Unmarshal", funcName)
		return nil, err
	}

	return &currency, nil
}
