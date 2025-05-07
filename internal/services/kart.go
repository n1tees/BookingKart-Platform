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
		Where("kartodrom_id = ? AND status = ?", kartodromID, models.Available).
		Find(&karts).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {

			return nil, errors.New("карты не найдены")

		} else {
			return nil, errors.New("ошибка при поиске картов")
		}
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
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("карт не найден")

			} else {
				return errors.New("ошибка при поиске карта")
			}
		}

		if kart.Status != models.Available {
			return errors.New("карт недоступен для бронирования")
		}

		if err := tx.Model(&kart).Update("status", models.InUse).Error; err != nil {
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
			if errors.Is(err, gorm.ErrRecordNotFound) {

				return errors.New("карт не найден")

			} else {
				return errors.New("ошибка при поиске карта")
			}
		}

		if err := tx.Model(&kart).Update("status", models.Available).Error; err != nil {
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
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("карт не найден")

			} else {
				return errors.New("ошибка при поиске карта")
			}
		}

		if err := tx.Model(&kart).Update("status", models.Broken).Error; err != nil {
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
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("карт не найден")

			} else {
				return errors.New("ошибка при поиске карта")
			}
		}

		if err := tx.Model(&kart).Update("status", models.Available).Error; err != nil {
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

			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("карт не найден")

			} else {
				return errors.New("ошибка при поиске карта")
			}
		}

		if err := tx.Model(&kart).Update("status", models.InStopList).Error; err != nil {
			return errors.New("ошибка при обновлении статуса карта")
		}

		return nil
	})
}
