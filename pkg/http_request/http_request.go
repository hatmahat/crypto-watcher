package http_request

import (
	"fmt"
	"io"
	"net/http"

	"github.com/sirupsen/logrus"
)

type Doer func(req *http.Request) (*http.Response, error)

func DoRequest(fn Doer, req *http.Request, funcName string) ([]byte, error) {
	resp, err := fn(req)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"err":       err.Error(),
			"resp_code": resp.StatusCode,
			"resp":      resp,
		}).Errorf("%s: Error Executing Request", funcName)
		return nil, err
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"err":  err.Error(),
			"resp": resp,
		}).Errorf("%s: Failed to Read Response", funcName)
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		logrus.WithFields(logrus.Fields{
			"resp_code": resp.StatusCode,
			"resp_body": string(responseBody),
		}).Errorf("%s: Non-OK HTTP status", funcName)
		return nil, fmt.Errorf("server response status: %d", resp.StatusCode)
	}

	return responseBody, nil
}
