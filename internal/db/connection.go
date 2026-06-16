package db

import (
	"fmt"

	"go-echo-starter/internal/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewGormDB(cfg config.DBConfig) (*gorm.DB, error) {
	// Use PostgreSQL driver instead of MySQL
	db, err := gorm.Open(postgres.Open(dsn(cfg)), &gorm.Config{
		Logger: newLoggerAdapter(),
	})
	if err != nil {
		return nil, fmt.Errorf("open db connection: %w", err)
	}

	return db, nil
}

func dsn(c config.DBConfig) string {
	// PostgreSQL DSN format
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		c.Host, c.User, c.Password, c.Name, c.Port,
	)
}
