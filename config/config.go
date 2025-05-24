package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	MySQLDSN string
	AppPort  string
}

func LoadConfig() (*AppConfig, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, errors.New("Failed to load .env file")
	}

	mysqlDSN := os.Getenv("MYSQL_DSN")
	if mysqlDSN == "" {
		mysqlDSN = "root:@tcp(127.0.0.1:3306)/movie_db?charset=utf8mb4&parseTime=True&loc=Local"
	}

	appPort := os.Getenv("APP_PORT")
	if appPort == "" {
		appPort = "8080"
	}

	return &AppConfig{
		MySQLDSN: mysqlDSN,
		AppPort: appPort,
	}, nil
}

func (c *AppConfig) GetDBDSN() string {
	return c.MySQLDSN
}