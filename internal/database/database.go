package database

import (
	"multi-chain-payment-gateway/internal/models"
	"strings"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Initialize(databaseURL string) (*gorm.DB, error) {
	// Extract database type and connection string
	var db *gorm.DB
	var err error

	if strings.HasPrefix(databaseURL, "sqlite://") {
		connStr := strings.TrimPrefix(databaseURL, "sqlite://")
		db, err = gorm.Open(sqlite.Open(connStr), &gorm.Config{})
	} else {
		// Default to SQLite
		db, err = gorm.Open(sqlite.Open("./payments.db"), &gorm.Config{})
	}

	if err != nil {
		return nil, err
	}

	// Auto-migrate schemas
	err = db.AutoMigrate(
		&models.Payment{},
		&models.PaymentOption{},
		&models.Transaction{},
	)
	if err != nil {
		return nil, err
	}

	return db, nil
}