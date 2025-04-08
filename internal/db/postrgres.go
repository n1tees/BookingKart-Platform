package db

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"BookingKart-Platform/config"
)

var DB *gorm.DB

func InitDB() {
	dsn := config.GetDBConnString()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Ошибка подключения к БД: %v", err)
	}

	DB = db
	log.Println("Успешное подключение к базе данных")
}
