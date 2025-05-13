package db

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"

	"github.com/n1tees/BookingKart-Platform/config"
	"github.com/n1tees/BookingKart-Platform/pkg/models"
)

var DB *gorm.DB

func InitDB() {
	dsn := config.GetDBConnString()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		log.Fatalf("Ошибка подключения к БД: %v", err)
	}
	log.Println("Успешное подключение к базе данных, подготовка к выполнению миграции")

	// if err := Migrate(db); err != nil {
	// 	log.Fatalf("Ошибка миграции базы данных: %v", err)
	// }

	DB = db
	log.Println("Пропуск миграций, не забыть вернуть")
	// log.Println("Успешное выполнение миграций")
}

func CloseDB() {
	sqlDB, err := DB.DB()
	if err != nil {
		log.Printf("Ошибка при получении SQL-соединения: %v", err)
		return
	}
	sqlDB.Close()
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
