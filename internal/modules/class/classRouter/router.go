package router

import (
	"github.com/gofiber/fiber/v2"
	classHandler "github.com/gunktp20/digital-hubx-be/internal/modules/class/classHandler"
)

func SetClassRoutes(api fiber.Router, classHttpHandler classHandler.ClassHttpHandlerService) {
	classRoute := api.Group("/class")

	classRoute.Get("/", classHttpHandler.GetAllClasses)
	classRoute.Get("/:class_id", classHttpHandler.GetClassById)

	// ? Admin Routes Group
	adminRoute := api.Group("/admin/class")
	adminRoute.Post("/", classHttpHandler.CreateClass)
	adminRoute.Put("/:class_id", classHttpHandler.UpdateClassDetails)
	adminRoute.Put("/:class_id/toggle-enable-question", classHttpHandler.ToggleClassEnableQuestion)
	adminRoute.Put("/:class_id/cover-image", classHttpHandler.UpdateClassCoverImage)
	adminRoute.Delete("/:class_id", classHttpHandler.DeleteClass)

}
