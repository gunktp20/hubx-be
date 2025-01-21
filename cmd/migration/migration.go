package main

import (
	"context"
	"log"

	"github.com/gunktp20/digital-hubx-be/pkg/config"
	"github.com/gunktp20/digital-hubx-be/pkg/constant"
	"github.com/gunktp20/digital-hubx-be/pkg/database"
	"github.com/gunktp20/digital-hubx-be/pkg/models"
)

func main() {
	ctx := context.Background()

	configPath := "../configuration"
	conf, err := config.GetConfig(configPath)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	db := database.NewGormPostgresDatabase(ctx, conf).Db

	if err := db.Exec("DROP TYPE IF EXISTS class_tiers CASCADE").Error; err != nil {
		log.Fatalln(constant.Red+"Failed to drop old enum class_tiers: "+constant.Reset, err)
	}

	if err := db.Exec("DROP TYPE IF EXISTS question_types CASCADE").Error; err != nil {
		log.Fatalln(constant.Red+"Failed to drop old enum question_types: "+constant.Reset, err)
	}

	if err := db.Exec("DROP TYPE IF EXISTS class_session_status CASCADE").Error; err != nil {
		log.Fatalln(constant.Red+"Failed to drop old enum class_session_status: "+constant.Reset, err)
	}

	if err := db.Exec("DROP TYPE IF EXISTS reg_status CASCADE").Error; err != nil {
		log.Fatalln(constant.Red+"Failed to drop old enum reg_status: "+constant.Reset, err)
	}

	if err := db.Exec("CREATE TYPE question_types AS ENUM ('choice', 'text')").Error; err != nil {
		log.Fatalln(constant.Red+"Failed to create enum question_types: "+constant.Reset, err)
	}

	if err := db.Exec("CREATE TYPE class_tiers AS ENUM ('essential', 'literacy', 'mastery')").Error; err != nil {
		log.Fatalln(constant.Red+"Failed to create enum class_tiers: "+constant.Reset, err)
	}

	if err := db.Exec("CREATE TYPE class_session_status AS ENUM ('available','closed','cancelled')").Error; err != nil {
		log.Fatalln(constant.Red+"Failed to create enum class_session_status: "+constant.Reset, err)
	}

	if err := db.Exec("CREATE TYPE reg_status AS ENUM ('registered','completed','cancelled')").Error; err != nil {
		log.Fatalln(constant.Red+"Failed to create enum reg_status: "+constant.Reset, err)
	}

	err = db.AutoMigrate(
		&models.UserSubQuestionAnswer{},
		&models.UserClassRegistration{},
		&models.UserQuestionAnswer{},
		&models.ClassHighLightImage{},

		&models.SubQuestionChoice{},
		&models.SubQuestion{},
		&models.Choice{},

		&models.Question{},

		&models.ClassSession{},
		&models.Class{},
		&models.Attendance{},
	)

	if err != nil {
		log.Fatalln(constant.Red+"Failed to migrate database: "+constant.Reset, err)
	} else {
		log.Fatalln(constant.Green + "Database migration successful" + constant.Reset)
	}

	defer func() {
		dbInstance, _ := db.DB()
		_ = dbInstance.Close()
	}()
}
