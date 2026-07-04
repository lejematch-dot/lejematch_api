package database

import (
	"Lejematch/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB() {
	var err error
	dsn := config.AppConfigInstance.DatabaseURL
	//dsn := "host=localhost user=postgres password=postgres dbname=postgres port=5432 sslmode=disable TimeZone=UTC"
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger:      logger.Default.LogMode(logger.Error),
		PrepareStmt: true,
	})
	if err != nil {
		panic("Failed to connect to database: " + err.Error())
	}
}

func CloseDB() {
	db, _ := DB.DB()
	db.Close()
}
