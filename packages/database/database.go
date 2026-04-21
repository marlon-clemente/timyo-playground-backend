package database

import (
	"fmt"

	"github.com/uptrace/opentelemetry-go-extra/otelgorm"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

func New(dsn string, development bool) (*gorm.DB, error) {
	if dsn == "" {
		return nil, fmt.Errorf("database: DSN is required (set DATABASE_DSN)")
	}

	logLevel := logger.Silent
	if development {
		logLevel = logger.Info
	}

	conn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		return nil, fmt.Errorf("database: failed to open connection: %w", err)
	}

	if err := conn.Use(otelgorm.NewPlugin()); err != nil {
		return nil, fmt.Errorf("database: failed to register otelgorm plugin: %w", err)
	}

	db = conn
	return db, nil
}

// Get returns the active database connection.
func Get() *gorm.DB {
	return db
}
