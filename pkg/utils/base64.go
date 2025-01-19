package utils

import (
	"encoding/base64"
	"fmt"
)

func DecodeBase64(content string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(content)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64 content: %w", err)
	}

	return string(data), nil
}

func DecodeBase64Byte(content string) ([]byte, error) {
	data, err := base64.StdEncoding.DecodeString(content)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64 content: %w", err)
	}
	return data, nil
}

func EncodeByteToBase64(content []byte) string {
	return base64.StdEncoding.EncodeToString(content)
}
