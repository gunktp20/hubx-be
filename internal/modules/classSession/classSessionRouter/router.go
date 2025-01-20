package router

import (
	"github.com/gofiber/fiber/v2"
	classSessionHandler "github.com/gunktp20/digital-hubx-be/internal/modules/classSession/classSessionHandler"
)

func SetClassSessionRoutes(api fiber.Router, classSessionHttpHandler classSessionHandler.ClassSessionHttpHandlerService) {

	classSessionRoute := api.Group("/class-session")

	classSessionRoute.Get("/", classSessionHttpHandler.GetAllClassSessions)

	// ? Admin Routes Group
	adminRoute := api.Group("/admin/class-session")
	adminRoute.Post("/", classSessionHttpHandler.CreateClassSession)
	adminRoute.Put("/:class_session_id/max-capacity", classSessionHttpHandler.SetMaxCapacity)
	adminRoute.Put("/:class_session_id/location", classSessionHttpHandler.UpdateClassSessionLocation)

}
