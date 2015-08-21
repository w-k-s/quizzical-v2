package utils

import (
	"net/http"
	"strconv"
)

type FormHelper struct {
	Request *http.Request
}

func (formHelper FormHelper) Int(key string, defaultValue int) int {

	value := formHelper.Request.FormValue(key)

	if len(value) == 0 {
		return defaultValue
	}

	result, err := strconv.Atoi(value)

	if err != nil {
		return defaultValue
	}

	return result
}

func (formHelper FormHelper) String(key string, defaultValue string) string {

	value := formHelper.Request.FormValue(key)

	if len(value) == 0 {
		return defaultValue
	}

	return value
}
