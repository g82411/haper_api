package models

import (
	"context"
	"fmt"
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

func NewDBClientWithContext(ctx context.Context) context.Context {
	con := ctx.Value("db")
	if con == nil {
		newConnection, err := NewDBClient()
		if err != nil {
			panic(fmt.Sprintf("NewDBClientWithContext: %v", err))
		}
		newContext := context.WithValue(ctx, "db", newConnection)
		return newContext
	}
	return ctx
}
