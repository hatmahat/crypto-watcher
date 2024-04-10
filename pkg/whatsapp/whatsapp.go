package whatsapp

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
	const funcName = "[pkg][whatsapp]SendWaMessageTemplate"

	request := MetaMessageRequest{
		MessaingProduct: whatsappConst,
		To:              phoneNumber,
		Type:            templateConst,
		Template: Template{
			Name: template,
			Language: Language{
				Code: enUsConst,
			},
		},
	}

	if len(parameters) != 0 {
		componentParameters := request.Template.Components[0].Parameters
		for _, param := range parameters {
			newParam := Parameter{
				Type: textConst,
				Text: param,
			}
			componentParameters = append(componentParameters, newParam)
		}
	}

	reqBody, err := json.Marshal(request)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"err":     err.Error(),
			"request": request,
		}).Error(funcName)
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf(sendWaMessageEndp, wm.phoneNumberId), bytes.NewBuffer(reqBody))
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"err":     err.Error(),
			"request": request,
		})
		return nil, err
	}

	auth := fmt.Sprintf("Bearer %s", wm.apiKey)
	req.Header.Set(http_const.ContentType, http_const.ApplicationJson)
	req.Header.Set(http_const.Authorization, auth)

	resp, err := wm.httpClient.Do(req)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"err":       err.Error(),
			"resp_code": resp.StatusCode,
			"resp":      resp,
		}).Error("Error Calling WhatsApp API")
		return nil, err
	}

	defer resp.Body.Close()
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"err":  err.Error(),
			"resp": resp,
		}).Error("Failed to Read Response")
		return nil, err
	}

	var metaMessageResponse MetaMessageResponse
	if err := json.Unmarshal(responseBody, &metaMessageResponse); err != nil {
		logrus.WithFields(logrus.Fields{
			"err":       err.Error(),
			"resp_body": string(responseBody),
		}).Error("Error Unmarshal to MetaMessageResponse")
		return nil, err
	}

	return &metaMessageResponse, nil
}
