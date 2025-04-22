package services

import (
	"github.com/n1tees/BookingKart-Platform/internal/db"
	"github.com/n1tees/BookingKart-Platform/internal/models"

	"errors"
)

// Получить список картодромов по городу
func GetKartodromsList(city *string) (*[]models.Kartodrom, error) {
	var kartodroms []models.Kartodrom

	query := db.DB.Preload("Schedules")

	if city != nil && *city != "" {
		query = query.Where("city = ?", *city)
	}

	if err := query.Find(&kartodroms).Error; err != nil {
		return nil, errors.New("ошибка при выводе списка картодромов")
	}

	return &kartodroms, nil
}
