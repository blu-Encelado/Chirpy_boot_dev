package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"net/http"
	"strings"
)

var ErrNoAuthHeaderIncluded = errors.New("no auth header included in request")
var ErrMalformedHeaderIncluded = errors.New("malformed header included in request")

func GetBearerToken(headers http.Header) (string, error) {
	token := headers.Get("Authorization")
	if token == "" {
		return "", ErrNoAuthHeaderIncluded
	}
	if !strings.HasPrefix(token, "Bearer ") {
		return "", ErrMalformedHeaderIncluded
	}
	string_list := strings.Split(token, " ")
	if len(string_list) < 2 {
		return "", ErrMalformedHeaderIncluded
	}
	string_var := string_list[1]
	return string_var, nil
}

func MakeRefreshToken() (string, error) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		return "", err
	}
	string_var := hex.EncodeToString(key)
	return string_var, nil
}
