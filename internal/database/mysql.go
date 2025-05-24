package database

import (
	"fmt"
	"log"
	"os"
	"roketin-case-study-challenge2/internal/entity"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitMySQLDB(dsn string) (*gorm.DB, error) {
	gormLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to MySQL: %w", err)
	}

	fmt.Println("Connected to MySQL database")

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("Failed to get SQL DB: %w", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	err = db.AutoMigrate(&entity.Movie{})
	if err != nil {
		return nil, fmt.Errorf("Failed to auto migrate: %w", err)
	}

	return db, nil
}
