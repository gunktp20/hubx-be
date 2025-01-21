// @title Digital HubX API
// @description Digital HubX API.
// @version 1.0.0
// @contact.name API Support (Tanapong R)
// @contact.url mailto:zTanapongR@pttep.com
// @contact.email zTanapongR@pttep.com
// @host localhost:3000
// @basePath /hubx-service

// @securityDefinitions.apiKey BearerAuth
// @in header
// @name Authorization
// @description This security definition is used for authenticating
package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gunktp20/digital-hubx-be/external/gcs"
	"github.com/gunktp20/digital-hubx-be/pkg/config"
	"github.com/gunktp20/digital-hubx-be/pkg/constant"
	"github.com/gunktp20/digital-hubx-be/pkg/database"
	"github.com/gunktp20/digital-hubx-be/server"
)

func main() {
	configPath := "../configuration"
	ctx := context.Background()

	conf := config.GetConfig(configPath)
	if conf == nil {
		log.Fatalln(constant.Red + "Failed to load configuration" + constant.Reset)
	}

	db := database.NewGormPostgresDatabase(ctx, conf)
	gcsClient := gcs.NewGcsClient(conf, nil)

	srv := server.NewFiberServer(conf, db, gcsClient)

	// Create a channel to listen for system signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	// Run the server in a separate goroutine
	go func() {
		srv.Start()
	}()

	log.Println("Server is running... Press Ctrl+C to shut down.")

	// Wait for termination signal
	<-quit
	log.Println("Gracefully shutting down...")

	// Shut down the server with a timeout context
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}

	log.Println("Server has been stopped successfully.")
}
