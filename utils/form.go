package utils

import (
	"net/http"
)

func FormValue(r *http.Request, key string, defaultValue string) string {

	if len(r.FormValue(key)) == 0 {
		return defaultValue
	}

	return r.FormValue(key)
}
