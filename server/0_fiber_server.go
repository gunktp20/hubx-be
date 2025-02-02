package server

import (
	"context"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"
	"github.com/gunktp20/digital-hubx-be/external/gcs"
	"github.com/gunktp20/digital-hubx-be/pkg/config"
	"github.com/gunktp20/digital-hubx-be/pkg/database"
	"github.com/gunktp20/digital-hubx-be/pkg/di"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	_ "github.com/gunktp20/digital-hubx-be/docs"
)

type (
	Server interface {
		Start()
		Shutdown(ctx context.Context) error
	}

	fiberServer struct {
		app       *fiber.App
		db        database.Database
		conf      *config.Config
		gcs       gcs.GcsClientService
		container *di.Container
	}
)

func NewFiberServer(conf *config.Config, db database.Database, gcs gcs.GcsClientService) Server {
	fiberApp := fiber.New(fiber.Config{
		ReadBufferSize:        60 * 1024,
		DisableStartupMessage: false,
	})

	container := di.NewContainer(conf, db)

	return &fiberServer{
		app:       fiberApp,
		db:        db,
		conf:      conf,
		gcs:       gcs,
		container: container,
	}
}

func (s *fiberServer) Start() {

	if s.conf.Logger.Enabled {
		s.app.Use(logger.New())
	}
	s.app.Use(cors.New(cors.Config{
		AllowOrigins:     s.conf.CORS.AllowOrigins,
		AllowMethods:     s.conf.CORS.AllowMethods,
		AllowHeaders:     s.conf.CORS.AllowHeaders,
		AllowCredentials: s.conf.CORS.AllowCredentials,
	}))
	s.app.Use(recover.New(recover.Config{EnableStackTrace: true}))

	// Default health check endpoint
	s.app.Get("", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Good health ✅",
		})
	})

	api := s.app.Group("/hubx-service")
	// Initialize routes (include Swagger)
	s.initializeRoutes(api)

	serverUrl := fmt.Sprintf(":%d", s.conf.Server.Port)
	if err := s.app.Listen(serverUrl); err != nil {
		panic(fmt.Sprintf("Failed to start server: %v", err))
	}
}

func (s *fiberServer) initializeRoutes(api fiber.Router) {

	if s.conf.Swagger.Enabled {
		s.app.Get("/swagger/*", swagger.HandlerDefault)
		log.Println("Swagger is enabled and available at /swagger/index.html")
	}

	// ? Ident middleware on top of http layer
	// s.app.Use(middleware.Ident) // TODO Remove this comment to enable bearer token validation for all entry endpoints.

	// ? initial http handler layer
	s.initializeClassHttpHandler(api, s.conf)
	s.initializeClassCategoryHttpHandler(api, s.conf)
	s.initializeClassSessionHttpHandler(api, s.conf)
	s.initializeClassRegistrationHttpHandler(api, s.conf)
	s.initializeQuestionHttpHandler(api, s.conf)
	s.initializeChoiceHttpHandler(api, s.conf)
	s.initializeUserQuestionAnswerHttpHandler(api, s.conf)
	s.initializeSubQuestionHttpHandler(api, s.conf)
	s.initializeSubQuestionChoiceHttpHandler(api, s.conf)
	s.initializeAttendanceHttpHandler(api, s.conf)
}

func (s *fiberServer) Shutdown(ctx context.Context) error {
	log.Println("Shutting down Fiber app...")

	if err := s.app.Shutdown(); err != nil {
		log.Printf("Error shutting down Fiber: %v", err)
		return err
	}

	log.Println("Closing database connection...")
	if err := s.db.Close(); err != nil {
		log.Printf("Error closing database: %v", err)
	}

	log.Println("Server shutdown complete.")
	return nil
}
