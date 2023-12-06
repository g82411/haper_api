package models

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"hyper_api/internal/config"
)

func NewDBClient() (*gorm.DB, error) {
	c := config.GetConfig()
	dbHost := c.DBHost
	dbName := c.DBName
	dbUsername := c.DBUsername
	dbPassword := c.DBPassword

	dsn := fmt.Sprint("host=", dbHost, " user=", dbUsername, " password=", dbPassword, " dbname=", dbName, " port=5432 sslmode=disable TimeZone=Asia/Taipei")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}
