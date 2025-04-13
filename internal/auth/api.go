package auth

import (
	"fmt"
	"net/http"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error) {
	h := headers.Get("Authorization")
	if h == "" {
		return "", fmt.Errorf("authorization header not found")
	}
	if !strings.HasPrefix(h, "ApiKey") {
		return "", fmt.Errorf("authorization header must start with ApiKey")
	}
	apiKey := strings.TrimPrefix(h, "ApiKey")
	apiKey = strings.TrimSpace(apiKey)

	return apiKey, nil
}
