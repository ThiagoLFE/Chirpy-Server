package auth

import (
	"fmt"
	"net/http"
	"strings"
)

func GetBearerToken(headers http.Header) (string, error) {
	authorization := headers.Get("Authorization")
	if authorization == "" {
		return "", fmt.Errorf("authorization is requiered")
	}
	token, bearer := strings.CutPrefix(authorization, "Bearer")
	if !bearer {
		return "", fmt.Errorf("bearer token not found")
	}

	return strings.TrimSpace(token), nil
}
