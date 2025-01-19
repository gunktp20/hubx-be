package router

import (
	"github.com/gofiber/fiber/v2"
	classSessionHandler "github.com/gunktp20/digital-hubx-be/internal/modules/classSession/classSessionHandler"
)

func SetClassSessionRoutes(api fiber.Router, classSessionHttpHandler classSessionHandler.ClassSessionHttpHandlerService) {
	routes := api.Group("/class-session")

	routes.Get("/", classSessionHttpHandler.GetAllClassSessions)
	routes.Post("/", classSessionHttpHandler.CreateClassSession)

}
