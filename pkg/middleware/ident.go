package middleware

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gunktp20/digital-hubx-be/pkg/constant"
)

const Bearer = "Bearer"

type JWTPayload struct {
	Name  string   `json:"name"`
	Email string   `json:"preferred_username"`
	Roles []string `json:"roles"`
}

func Ident(c *fiber.Ctx) error {
	auth := strings.Fields(c.Get(fiber.HeaderAuthorization))
	if len(auth) != 2 || !strings.EqualFold(auth[0], Bearer) {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	tokens := strings.Split(auth[1], ".")
	if len(tokens) < 2 {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	encoded := tokens[1]
	if l := len(encoded) % 4; l > 0 {
		encoded += strings.Repeat("=", 4-l)
	}
	decoded, err := base64.URLEncoding.DecodeString(encoded)
	if err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	var payload JWTPayload
	if err := json.Unmarshal(decoded, &payload); err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	ctx := context.WithValue(c.UserContext(), constant.CtxToken, auth[1])
	ctx = context.WithValue(ctx, constant.CtxEmail, payload.Email)
	ctx = context.WithValue(ctx, constant.CtxName, payload.Name)
	ctx = context.WithValue(ctx, constant.CtxRoles, payload.Roles)
	c.SetUserContext(ctx)
	fmt.Print("================================================================================================================================================")
	fmt.Print(payload.Email)
	c.Locals("token", auth[1])
	c.Locals("email", payload.Email)
	c.Locals("name", payload.Name)
	c.Locals("roles", payload.Roles)
	return c.Next()
}
