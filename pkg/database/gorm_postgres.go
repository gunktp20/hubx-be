package database

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gunktp20/digital-hubx-be/pkg/config"
	"github.com/gunktp20/digital-hubx-be/pkg/constant"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type gormPostgresDatabase struct {
	Db *gorm.DB
}

var (
	once       sync.Once
	dbInstance *gormPostgresDatabase
)

func NewGormPostgresDatabase(pctx context.Context, conf *config.Config) *gormPostgresDatabase {
	ctx, cancel := context.WithTimeout(pctx, 20*time.Second)
	defer cancel()

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		conf.Db.DbHost,
		conf.Db.DbPort,
		conf.Db.DbUser,
		conf.Db.DbSecret,
		conf.Db.DbName,
		conf.Db.SSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // TODO Uncomment for monitor SQL query
	})
	if err != nil {
		log.Fatalln(constant.Red+"Gorm Postgres database connection failed"+constant.Reset, err)
	}

	postgresDB, err := db.DB()
	if err != nil {
		log.Fatalln(constant.Red+"Gorm Postgres database get sqlDB instance failed"+constant.Reset, err)
	}

	if err := postgresDB.PingContext(ctx); err != nil {
		log.Fatalln(constant.Red+"Gorm Postgres database ping failed"+constant.Reset, err)
	}

	dbInstance = &gormPostgresDatabase{Db: db}

	fmt.Println(constant.Green + "Postgres database connection successful" + constant.Reset)
	return dbInstance
}

func (p *gormPostgresDatabase) GetDb() *gorm.DB {
	return dbInstance.Db
}

func (p *gormPostgresDatabase) Close() error {
	postgresDB, err := p.Db.DB()
	if err != nil {
		return err // กรณีที่ไม่สามารถดึง instance ได้
	}
	return postgresDB.Close() // ปิดการเชื่อมต่อฐานข้อมูล
}
