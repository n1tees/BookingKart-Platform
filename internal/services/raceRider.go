package services

import (
	"errors"

	"github.com/n1tees/BookingKart-Platform/internal/db"
	"github.com/n1tees/BookingKart-Platform/internal/models"
	"gorm.io/gorm"
)

// Добавить райдера в заезд
func RegisterRider(raceID, riderID uint) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		// Проверяем гонку
		var race models.Race
		if err := tx.First(&race, raceID).Error; err != nil {
			return errors.New("гонка не найдена")
		}

		// Проверяем пользователя
		var rider models.User
		if err := tx.First(&rider, riderID).Error; err != nil {
			return errors.New("пользователь не найден")
		}

		// Проверяем, не зарегистрирован ли уже райдер
		var existing models.RaceRider
		if err := tx.Where("race_id = ? AND rider_id = ?", raceID, riderID).First(&existing).Error; err == nil {
			return errors.New("пользователь уже зарегистрирован в этой гонке")
		}

		// Регистрируем райдера
		newRider := models.RaceRider{
			RaceID:         raceID,
			RiderID:        riderID,
			ResultTypeID:   1, // по умолчанию "Участник", ID можно будет уточнить
			PersonalResult: 0,
		}

		if err := tx.Create(&newRider).Error; err != nil {
			return errors.New("ошибка при регистрации участника")
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
			return errors.New("участник не найден в этой гонке")
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
			return errors.New("участник не найден в этой гонке")
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
