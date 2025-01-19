package router

import (
	"github.com/gofiber/fiber/v2"
	classRegistrationHandler "github.com/gunktp20/digital-hubx-be/internal/modules/classRegistration/classRegistrationHandler"
)

func SetClassRegistrationRoutes(api fiber.Router, classRegistrationHttpHandler classRegistrationHandler.ClassRegistrationHttpHandlerService) {
	routes := api.Group("/class-registration")

	routes.Get("/", classRegistrationHttpHandler.GetUserRegistrations)
	routes.Post("/", classRegistrationHttpHandler.CreateClassRegistration)

}
