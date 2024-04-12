package whatsapp_cloud_api

import (
	"bytes"
	"context"
	"crypto-watcher-backend/internal/constant/http_const"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/sirupsen/logrus"
)

type (
	WaMessaging interface {
		SendWaMessageByTemplate(ctx context.Context, phoneNumber, template string, parameters []string) (*MetaMessageResponse, error)
	}

	waMessaging struct {
		host          string
		apiKey        string
		phoneNumberId string
		httpClient    *http.Client
	}
)

func NewWaMessaging(host, apiKey, phoneNumberId string) WaMessaging {
	httpClient := &http.Client{}

	return &waMessaging{
		host:          host,
		apiKey:        apiKey,
		phoneNumberId: phoneNumberId,
		httpClient:    httpClient,
	}
}

func (wm *waMessaging) SendWaMessageByTemplate(ctx context.Context, phoneNumber, template string, parameters []string) (*MetaMessageResponse, error) {
	const funcName = "[pkg][whatsapp_cloud_api]SendWaMessageTemplate"

	request := MetaMessageRequest{
		MessaingProduct: whatsappConst,
		To:              phoneNumber,
		Type:            templateConst,
		Template: Template{
			Name: template,
			Language: Language{
				Code: enUsConst,
			},
			Components: []Component{
				{
					Type: bodyConst,
				},
			},
		},
	}

	componentParameters := make([]Parameter, 0)
	if len(parameters) != 0 {
		for _, param := range parameters {
			newParam := Parameter{
				Type: textConst,
				Text: param,
			}
			componentParameters = append(componentParameters, newParam)
		}
	}

	request.Components[0].Parameters = componentParameters

	reqBody, err := json.Marshal(request)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"err":     err.Error(),
			"request": request,
		}).Error(funcName)
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, wm.host+fmt.Sprintf(sendWaMessagePath, wm.phoneNumberId), bytes.NewBuffer(reqBody))
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"err":     err.Error(),
			"request": request,
		}).Errorf("Error Making Request: %s", funcName)
		return nil, err
	}

	req.Header.Set(http_const.ContentType, http_const.ApplicationJson)
	req.Header.Set(http_const.Authorization, fmt.Sprintf(http_const.Bearer, wm.apiKey))

	resp, err := wm.httpClient.Do(req)
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

	var metaMessageResponse MetaMessageResponse
	if err := json.Unmarshal(responseBody, &metaMessageResponse); err != nil {
		logrus.WithFields(logrus.Fields{
			"err":       err.Error(),
			"resp_body": string(responseBody),
		}).Errorf("Error Unmarshal: %s", funcName)
		return nil, err
	}

	return &metaMessageResponse, nil
}
