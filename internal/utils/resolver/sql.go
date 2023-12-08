package resolver

import (
	"fmt"
	"hyper_api/internal/config"
)

func FormatSQLDSNFromConfig() string {
	env := config.GetConfig()
	dbHost := env.DBHost
	dbName := env.DBName
	dbUsername := env.DBUsername
	dbPassword := env.DBPassword
	dsn := fmt.Sprint("host=", dbHost, " user=", dbUsername, " password=", dbPassword, " dbname=", dbName, " port=5432 sslmode=disable TimeZone=Asia/Taipei")
	return dsn
}
