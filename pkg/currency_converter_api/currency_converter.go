package currency_converter_api

import (
	"context"
	"crypto-watcher-backend/pkg/http_request"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"

	"github.com/sirupsen/logrus"
)

type (
	CurrencyConverter interface {
		GetCurrencyRate(ctx context.Context, params map[string]string) (*CurrencyConverterResponse, error)
	}

	currencyConverter struct {
		host       string
		apiKey     string
		httpClient *http.Client
	}
)

func NewCurrencyConverter(host, apiKey string) CurrencyConverter {
	httpClient := &http.Client{}

	return &currencyConverter{
		host:       host,
		apiKey:     apiKey,
		httpClient: httpClient,
	}
}

func (cc *currencyConverter) GetCurrencyRate(ctx context.Context, queryParams map[string]string) (*CurrencyConverterResponse, error) {
	const funcName = "[pkg][currency_converter_api]GetCurrencyRate"

	u, err := url.Parse(cc.host)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"err": err.Error(),
		}).Errorf("%s: Error Parsing URL", funcName)
		return nil, err
	}

	u.Path += currencyConverterPath

	q := u.Query()
	for key, value := range queryParams {
		q.Set(key, value)
	}
	u.RawQuery = q.Encode()

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"err": err.Error(),
		}).Errorf("%s: Error Making Request", funcName)
		return nil, err
	}

	cleanedUrl := strings.TrimPrefix(cc.host, "https://")
	req.Header.Set(XRapidAPIHost, cleanedUrl)
	req.Header.Set(XRapidAPIKey, cc.apiKey)

	responseBody, err := http_request.DoRequest(cc.httpClient.Do, req, funcName)
	if err != nil {
		logrus.WithError(err).Errorf("%s: Error Do Request", funcName)
		return nil, err
	}

	var currencyConverterResponse CurrencyConverterResponse
	if err := json.Unmarshal(responseBody, &currencyConverterResponse); err != nil {
		logrus.WithFields(logrus.Fields{
			"err":       err.Error(),
			"resp_body": string(responseBody),
		}).Errorf("%s: Error Unmarshal", funcName)
		return nil, err
	}

	return &currencyConverterResponse, nil
}
