package middleware

import (
	"log"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gunktp20/digital-hubx-be/pkg/response"
	"github.com/gunktp20/digital-hubx-be/pkg/utils"
)

func DecodeJwtPayload(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return response.ErrResponse(c, http.StatusUnauthorized, "Authorization header is missing", nil)
	}

	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		return response.ErrResponse(c, http.StatusUnauthorized, "Invalid authorization format", nil)
	}

	claims, err := utils.ParseJwt(tokenParts[1])
	if err != nil {
		return response.ErrResponse(c, http.StatusUnauthorized, err.Error(), nil)
	}

	// ดึง user_id จาก claims และเก็บใน context
	if email, ok := claims["upn"].(string); ok {
		c.Locals("email", email)
		log.Printf("email : %s", email)
	} else {
		return response.ErrResponse(c, http.StatusUnauthorized, "email not found in token claims", nil)
	}

	return c.Next()

}
