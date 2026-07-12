package auth

import (
	"fmt"
	"net/http"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error) {
	authorization := headers.Get("Authorization")
	if authorization == "" {
		return "", fmt.Errorf("Authorization is required")
	}
	apiKey, ok := strings.CutPrefix(authorization, "ApiKey")
	if !ok {
		return "", fmt.Errorf("Invalid authorization key")
	}

	return strings.TrimSpace(apiKey), nil
}
