package db

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/n1tees/BookingKart-Platform/config"
	"github.com/n1tees/BookingKart-Platform/internal/models"
)

var DB *gorm.DB

func InitDB() {
	dsn := config.GetDBConnString()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Ошибка подключения к БД: %v", err)
	}

	if err := Migrate(db); err != nil {
		log.Fatalf("Ошибка миграции базы данных: %v", err)
	}

	DB = db
	log.Println("Успешное подключение к базе данных и выполнены миграции")
}

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.AuthCredential{},
		&models.User{},
		&models.Profile{},
		&models.Payment{},
		&models.RaceRider{},
		&models.Booking{},
		&models.KartBooking{},
		&models.ResultType{},
		&models.Race{},
		&models.Track{},
		&models.Kart{},
		&models.TrackStat{},
		&models.KartModel{},
		&models.CommonStat{},
		&models.KartodromSchedule{},
		&models.Kartodrom{},
		&models.KartStat{},
	)
}
