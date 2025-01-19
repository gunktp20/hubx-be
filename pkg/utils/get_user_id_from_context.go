package utils

import (
	"errors"

	"github.com/gofiber/fiber/v2"
)

func GetUserEmailFromContext(c *fiber.Ctx) (string, error) {
	userEmail := c.Locals("email")
	if userEmail == nil {
		return "", errors.New("email not found in context")
	}

	uEmailStr, ok := userEmail.(string)
	if !ok {
		return "", errors.New("email is not a valid string")
	}

	return uEmailStr, nil
}
