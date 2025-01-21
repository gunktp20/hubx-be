package router

import (
	"github.com/gofiber/fiber/v2"
	classSessionHandler "github.com/gunktp20/digital-hubx-be/internal/modules/classSession/classSessionHandler"
	"github.com/gunktp20/digital-hubx-be/pkg/middleware"
)

func SetClassSessionRoutes(api fiber.Router, classSessionHttpHandler classSessionHandler.ClassSessionHttpHandlerService) {

	classSessionRoute := api.Group("/class-session")

	classSessionRoute.Get("/", classSessionHttpHandler.GetAllClassSessions)

	// ? Admin Routes
	adminRoute := api.Group("/admin/class-session", middleware.Ident, middleware.PermissionCheck)
	adminRoute.Post("/", classSessionHttpHandler.CreateClassSession)
	adminRoute.Delete("/:class_session_id", classSessionHttpHandler.DeleteClassSessionByID)
	adminRoute.Put("/:class_session_id/max-capacity", classSessionHttpHandler.SetMaxCapacity)
	adminRoute.Put("/:class_session_id/location", classSessionHttpHandler.UpdateClassSessionLocation)

}
