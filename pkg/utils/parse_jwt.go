package utils

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
)

func ParseJwt(tokenString string) (map[string]interface{}, error) {

	// แยก JWT ออกเป็นส่วนๆ (Header, Payload, Signature)
	parts := strings.Split(tokenString, ".")
	if len(parts) < 2 {
		return nil, errors.New("invalid token format")
	}

	// Decode ส่วน Payload (ส่วนที่สองของ JWT)
	payload := parts[1]

	// แก้ไข Base64 ให้ถูกต้อง (เพิ่ม padding ถ้าจำเป็น)
	if mod := len(payload) % 4; mod > 0 {
		payload += strings.Repeat("=", 4-mod)
	}

	// ถอดรหัส Base64
	decoded, err := base64.URLEncoding.DecodeString(payload)
	if err != nil {
		return nil, errors.New("failed to decode payload")
	}

	// แปลง JSON Payload เป็น map[string]interface{}
	var claims map[string]interface{}
	if err := json.Unmarshal(decoded, &claims); err != nil {
		return nil, errors.New("failed to parse payload as JSON")
	}

	return claims, nil
}
