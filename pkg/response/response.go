package response

import "github.com/gofiber/fiber/v2"

type (
	MsgResponse struct {
		Message string             `json:"message" example:"Operation successful"`                    // The response message
		Status  int                `json:"status" example:"200"`                                      // HTTP status code
		Details *map[string]string `json:"details,omitempty" example:"{\"field\": \"error detail\"}"` // Additional details (optional)
	}
)

func ErrResponse(c *fiber.Ctx, statusCode int, message string, details *map[string]string) error {
	return c.Status(statusCode).JSON(&MsgResponse{Status: statusCode, Message: message, Details: details})
}

func SuccessResponse(c *fiber.Ctx, statusCode int, data any) error {
	return c.Status(statusCode).JSON(&MsgResponse{Status: statusCode, Message: "Operation successful"})
}
