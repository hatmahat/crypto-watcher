package validation

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

func ValidateFromMapper(key string, mapper map[string]string) (*string, error) {
	const funcName = "[pkg][validation]ValidateFromMapper"

	value, ok := mapper[key]
	if !ok {
		logrus.Errorf("%s: Key [%s] Not Found", funcName, key)
		return nil, fmt.Errorf("%s: Key [%s] Not Found", funcName, key)
	}
	return &value, nil
}
