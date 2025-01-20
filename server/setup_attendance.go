package server

import (
	"github.com/gofiber/fiber/v2"
	attendanceHandler "github.com/gunktp20/digital-hubx-be/internal/modules/attendance/attendanceHandler"
	attendanceRouter "github.com/gunktp20/digital-hubx-be/internal/modules/attendance/attendanceRouter"
	attendanceUsecase "github.com/gunktp20/digital-hubx-be/internal/modules/attendance/attendanceUsecase"

	"github.com/gunktp20/digital-hubx-be/pkg/config"
)

func (s *fiberServer) initializeAttendanceHttpHandler(api fiber.Router, conf *config.Config) {
	// ? Initialize all layers

	attendanceUsecase := attendanceUsecase.NewAttendanceUsecase(
		s.container.Repositories.AttendanceRepo,
		s.container.Repositories.ClassSessionRepo,
		s.container.Repositories.ClassRegistrationRepo,
	)
	attendanceHttpHandler := attendanceHandler.NewAttendanceHttpHandler(attendanceUsecase)

	// Routers
	attendanceRouter.SetAttendanceRoute(api, attendanceHttpHandler)
}
