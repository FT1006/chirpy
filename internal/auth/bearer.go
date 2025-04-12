package auth

import (
	"fmt"
	"net/http"
	"strings"
)

func GetBearerToken(headers http.Header) (string, error) {
	h := headers.Get("Authorization")
	if h == "" {
		return "", fmt.Errorf("authorization header not found")
	}
	if !strings.HasPrefix(h, "Bearer ") {
		return "", fmt.Errorf("authorization header must start with Bearer")
	}
	tokenSharing := strings.TrimPrefix(h, "Bearer")
	tokenSharing = strings.TrimSpace(tokenSharing)

	return tokenSharing, nil
}
