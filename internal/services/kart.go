package services

import (
	"errors"
	"sync"

	"github.com/n1tees/BookingKart-Platform/internal/db"
	"github.com/n1tees/BookingKart-Platform/internal/models"

	"gorm.io/gorm"
)

// Блокировка для операций с картами
var kartMu sync.Mutex

// Получить список доступных картов в определенном картодроме
func GetAvailableKarts(kartodromID uint) (*[]models.Kart, error) {
	var karts []models.Kart

	if err := db.DB.
		Where("kartodrom_id = ? AND kart_status = ?", kartodromID, models.Availible).
		Find(&karts).Error; err != nil {
		return nil, errors.New("ошибка при получении доступных картов")
	}

	return &karts, nil
}

// Бронирование карта — меняем статус на InUse
func BookKart(kartID uint) error {
	kartMu.Lock()
	defer kartMu.Unlock()

	return db.DB.Transaction(func(tx *gorm.DB) error {
		var kart models.Kart
		if err := tx.First(&kart, kartID).Error; err != nil {
			return errors.New("карт не найден")
		}

		if kart.Status != models.Availible {
			return errors.New("карт недоступен для бронирования")
		}

		if err := tx.Model(&kart).Update("kart_status", models.InUse).Error; err != nil {
			return errors.New("ошибка при обновлении статуса карта")
		}

		return nil
	})
}

// Освобождение карта — меняем статус на Availible
func FreeKart(kartID uint) error {
	kartMu.Lock()
	defer kartMu.Unlock()

	return db.DB.Transaction(func(tx *gorm.DB) error {
		var kart models.Kart
		if err := tx.First(&kart, kartID).Error; err != nil {
			return errors.New("карт не найден")
		}

		if err := tx.Model(&kart).Update("kart_status", models.Availible).Error; err != nil {
			return errors.New("ошибка при обновлении статуса карта")
		}

		return nil
	})
}

// Помечаем карт как сломанный — статус Broken
func SetKartBroken(kartID uint) error {
	kartMu.Lock()
	defer kartMu.Unlock()

	return db.DB.Transaction(func(tx *gorm.DB) error {
		var kart models.Kart
		if err := tx.First(&kart, kartID).Error; err != nil {
			return errors.New("карт не найден")
		}

		if err := tx.Model(&kart).Update("kart_status", models.Broken).Error; err != nil {
			return errors.New("ошибка при обновлении статуса карта")
		}

		return nil
	})
}

// Восстанавливаем карт после ремонта — статус Availible
func RepairKart(kartID uint) error {
	kartMu.Lock()
	defer kartMu.Unlock()

	return db.DB.Transaction(func(tx *gorm.DB) error {
		var kart models.Kart
		if err := tx.First(&kart, kartID).Error; err != nil {
			return errors.New("карт не найден")
		}

		if err := tx.Model(&kart).Update("kart_status", models.Availible).Error; err != nil {
			return errors.New("ошибка при обновлении статуса карта")
		}

		return nil
	})
}

// Помещаем карт в стоп-лист — статус InStopList
func SetKartInStopList(kartID uint) error {
	kartMu.Lock()
	defer kartMu.Unlock()

	return db.DB.Transaction(func(tx *gorm.DB) error {
		var kart models.Kart
		if err := tx.First(&kart, kartID).Error; err != nil {
			return errors.New("карт не найден")
		}

		if err := tx.Model(&kart).Update("kart_status", models.InStopList).Error; err != nil {
			return errors.New("ошибка при обновлении статуса карта")
		}

		return nil
	})
}
