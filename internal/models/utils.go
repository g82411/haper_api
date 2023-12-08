package models

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"hyper_api/internal/utils/resolver"
)

func NewDBClient() (*gorm.DB, error) {
	dsn := resolver.FormatSQLDSNFromConfig()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}
