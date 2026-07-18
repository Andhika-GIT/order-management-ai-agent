package database

import (
	"fmt"

	"github.com/Andhika-GIT/go-message-broker-monorepo/configs"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDatabase(cfg *configs.DatabaseConfig) *gorm.DB {
	DB_HOST := cfg.Host
	DB_NAME := cfg.Name
	DB_PORT := cfg.Port
	DB_USERNAME := cfg.User
	DB_PASSWORD := cfg.Password

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Shanghai", DB_HOST, DB_USERNAME, DB_PASSWORD, DB_NAME, DB_PORT)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic(fmt.Errorf("fatal error connecting to database: %w", err))
	}

	return db

}
