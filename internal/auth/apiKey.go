package auth

import (
	"net/http"
	"strings"
)

func GetApiKey(headers http.Header) (string, error) {
	apiKey := headers.Get("Authorization")
	if apiKey == "" {
		return "", ErrNoAuthHeaderIncluded
	}
	if !strings.HasPrefix(apiKey, "ApiKey ") {
		return "", ErrMalformedHeaderIncluded
	}
	string_list := strings.Split(apiKey, " ")
	if len(string_list) < 2 {
		return "", ErrMalformedHeaderIncluded
	}
	string_var := string_list[1]
	return string_var, nil
}
