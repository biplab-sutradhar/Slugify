package auth

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

func GenerateAPIKey() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %v", err)
	}
	return base64.RawURLEncoding.EncodeToString(bytes), nil
}

func ValidateAPIKey(key string) error {
	if len(key) != 43 {
		return fmt.Errorf("invalid API key length")
	}
	return nil
}
