package middleware

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	userRepository "github.com/gunktp20/digital-hubx-be/internal/modules/user/userRepository"
	"github.com/gunktp20/digital-hubx-be/pkg/response"
)

func EnsureUserMiddleware(userRepo userRepository.UserRepositoryService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		email := c.Locals("email").(string)

		user, err := userRepo.GetOrCreateUser(email)
		if err != nil {
			return response.ErrResponse(c, http.StatusInternalServerError, err.Error(), nil)
		}

		c.Locals("user", user)
		return c.Next()
	}
}
