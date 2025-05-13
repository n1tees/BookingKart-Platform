package services

import (
	"errors"
	"fmt"

	"github.com/n1tees/BookingKart-Platform/pkg/db"
	"github.com/n1tees/BookingKart-Platform/pkg/models"
	"gorm.io/gorm"
)

// Добавить райдера в заезд
// Добавить райдера в заезд с указанием типа результата
func RegisterRider(raceID, riderID, resultTypeID uint) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		var race models.Race
		if err := tx.First(&race, raceID).Error; err != nil {
			return fmt.Errorf("гонка не найдена: %w", err)
		}

		var rider models.User
		if err := tx.First(&rider, riderID).Error; err != nil {
			return fmt.Errorf("пользователь не найден: %w", err)
		}

		var existing models.RaceRider
		if err := tx.Where("race_id = ? AND rider_id = ?", raceID, riderID).First(&existing).Error; err == nil {
			return errors.New("пользователь уже зарегистрирован в гонке")
		}

		newRider := models.RaceRider{
			RaceID:         raceID,
			RiderID:        riderID,
			ResultTypeID:   resultTypeID, // <---- исправление
			PersonalResult: 0,
		}

		if err := tx.Create(&newRider).Error; err != nil {
			return fmt.Errorf("ошибка при регистрации участника: %w", err)
		}

		return nil
	})
}

// Удалить из заезда
func RemoveRider(raceID, riderID uint) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		// Проверяем существование райдера в гонке
		var rider models.RaceRider
		if err := tx.Where("race_id = ? AND rider_id = ?", raceID, riderID).First(&rider).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("пользователь не найден")
			}
			return errors.New("ошибка при поиске пользователя")
		}

		// Удаляем райдера
		if err := tx.Delete(&rider).Error; err != nil {
			return errors.New("ошибка при удалении участника")
		}

		return nil
	})
}

// Добавление результата
func AddRaceResult(raceID, riderID, resultTypeID, personalResult uint) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		// Находим участника в гонке
		var rider models.RaceRider
		if err := tx.Where("race_id = ? AND rider_id = ?", raceID, riderID).First(&rider).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("пользователь не найден")
			}
			return errors.New("ошибка при поиске пользователя")
		}

		// Обновляем результат
		rider.ResultTypeID = resultTypeID
		rider.PersonalResult = personalResult

		if err := tx.Save(&rider).Error; err != nil {
			return errors.New("ошибка при сохранении результата участника")
		}

		return nil
	})
}
