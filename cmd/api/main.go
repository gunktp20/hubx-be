// @title Digital HubX API
// @description Digital HubX API.
// @version 1.0.0
// @contact.name API Support (Tanapong R)
// @contact.url mailto:zTanapongR@pttep.com
// @contact.email zTanapongR@pttep.com
// @host hrhyperspace-dev.api.pttep.com
// @basePath /hubx-service

// @securityDefinitions.apiKey	BearerAuth
// @in							header
// @name						Authorization
// @description		This security definition is used for authenticating API requests using a Bearer token. Clients must include the token in the Authorization header formatted as "Bearer {token}".
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

	// Channel สำหรับจับสัญญาณ
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	// รันเซิร์ฟเวอร์ใน Goroutine
	go func() {
		srv.Start()
	}()

	log.Println("Server is running... Press Ctrl+C to shut down.")

	// รอจนกว่าจะมีสัญญาณ
	<-quit
	log.Println("Gracefully shutting down...")

	// ทำการปิดเซิร์ฟเวอร์
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}

	log.Println("Server has been stopped successfully.")
}
