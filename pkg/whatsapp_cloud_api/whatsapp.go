package whatsapp_cloud_api

import (
	"bytes"
	"context"
	"crypto-watcher-backend/internal/constant/http_const"
	"crypto-watcher-backend/pkg/http_request"
	"encoding/json"
	"fmt"
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
		}).Errorf("%s: Error Marshalling", funcName)
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, wm.host+fmt.Sprintf(sendWaMessagePath, wm.phoneNumberId), bytes.NewBuffer(reqBody))
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"err":     err.Error(),
			"request": request,
		}).Errorf("%s: Error Making Request", funcName)
		return nil, err
	}

	req.Header.Set(http_const.ContentType, http_const.ApplicationJson)
	req.Header.Set(http_const.Authorization, fmt.Sprintf(http_const.Bearer, wm.apiKey))

	responseBody, err := http_request.DoRequest(wm.httpClient.Do, req, funcName)
	if err != nil {
		logrus.WithError(err).Errorf("%s: Error Do Request", funcName)
		return nil, err
	}

	var metaMessageResponse MetaMessageResponse
	if err := json.Unmarshal(responseBody, &metaMessageResponse); err != nil {
		logrus.WithFields(logrus.Fields{
			"err":       err.Error(),
			"resp_body": string(responseBody),
		}).Errorf("%s: Error Unmarshal", funcName)
		return nil, err
	}

	return &metaMessageResponse, nil
}
