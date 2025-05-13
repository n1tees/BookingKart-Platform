package services

import (
	"errors"
	"time"

	"github.com/n1tees/BookingKart-Platform/internal/db"
	"github.com/n1tees/BookingKart-Platform/internal/models"
	"gorm.io/gorm"
)

// Получить достпуные для бронирования карты
func GetAvailableKartsForBooking(kartodromID uint, start time.Time, end time.Time) (*[]models.Kart, error) {
	// 1. Сначала находим все карты в нужном картодроме, которые не сломаны и не в стоп-листе
	var allKarts []models.Kart
	if err := db.DB.Where("kartodrom_id = ? AND status = ?", kartodromID, models.Available).
		Find(&allKarts).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("карт не найден")

		} else {
			return nil, errors.New("ошибка при поиске карта")
		}
	}

	if len(allKarts) == 0 {
		return &[]models.Kart{}, nil
	}

	// Собираем IDs всех найденных картов
	var kartIDs []uint
	for _, kart := range allKarts {
		kartIDs = append(kartIDs, kart.ID)
	}

	// 2. Ищем карты, которые уже забронированы в это время
	var busyKartIDs []uint
	if err := db.DB.Model(&models.KartBooking{}).
		Joins("JOIN bookings ON bookings.id = kart_bookings.booking_id").
		Where("kart_bookings.kart_id IN ? AND bookings.status IN ? AND bookings.start_time < ? AND bookings.end_time > ?",
			kartIDs, []models.BookingStatus{models.BookingActive, models.BookingReserve}, end, start).
		Pluck("kart_bookings.kart_id", &busyKartIDs).Error; err != nil {
		return nil, errors.New("ошибка при проверке занятых картов")
	}

	// 3. Убираем занятые карты из общего списка
	var availableKarts []models.Kart
	for _, kart := range allKarts {
		isBusy := false
		for _, busyID := range busyKartIDs {
			if kart.ID == busyID {
				isBusy = true
				break
			}
		}
		if !isBusy {
			availableKarts = append(availableKarts, kart)
		}
	}

	return &availableKarts, nil
}

// Бронь карта
func ReserveKart(bookingID uint, kartID uint) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		// Проверяем существование карты
		var kart models.Kart
		if err := tx.First(&kart, kartID).Error; err != nil {
			return errors.New("карт не найден")
		}

		// Привязываем карт к бронированию
		kartBooking := models.KartBooking{
			BookingID: bookingID,
			KartID:    kartID,
		}

		if err := tx.Create(&kartBooking).Error; err != nil {
			return errors.New("ошибка при резервировании карта")
		}

		// Обновляем статус карты
		if err := tx.Model(&kart).Update("kart_status", models.InUse).Error; err != nil {
			return errors.New("ошибка при обновлении статуса карта")
		}

		return nil
	})
}

// Активация брони карта
func ActivateKartBooking(bookingID uint, kartID uint) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		var booking models.Booking
		if err := tx.First(&booking, bookingID).Error; err != nil {
			return errors.New("бронирование не найдено")
		}

		if booking.Status != models.BookingReserve {
			return errors.New("только резервированное бронирование можно активировать")
		}

		return tx.Model(&booking).Update("status", models.BookingActive).Error
	})
}

// Завершение брони карт
func FinishKartBooking(bookingID uint, kartID uint) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		var booking models.Booking
		if err := tx.First(&booking, bookingID).Error; err != nil {
			return errors.New("бронирование не найдено")
		}

		if booking.Status != models.BookingActive {
			return errors.New("только активное бронирование можно завершить")
		}

		// Завершаем бронирование
		if err := tx.Model(&booking).Update("status", models.BookingClose).Error; err != nil {
			return errors.New("ошибка при завершении бронирования")
		}

		// Освобождаем карт
		if err := tx.Model(&models.Kart{}).Where("id = ?", kartID).Update("kart_status", models.Available).Error; err != nil {
			return errors.New("ошибка при освобождении карта")
		}

		return nil
	})
}

// Отмена брони карта
func CancelKartBooking(bookingID uint, kartID uint) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		var booking models.Booking
		if err := tx.First(&booking, bookingID).Error; err != nil {
			return errors.New("бронирование не найдено")
		}

		if booking.Status != models.BookingActive && booking.Status != models.BookingReserve {
			return errors.New("нельзя отменить завершённое бронирование")
		}

		// Отменяем бронирование
		if err := tx.Model(&booking).Update("status", models.BookingCancel).Error; err != nil {
			return errors.New("ошибка при отмене бронирования")
		}

		// Освобождаем карт
		if err := tx.Model(&models.Kart{}).Where("id = ?", kartID).Update("kart_status", models.Available).Error; err != nil {
			return errors.New("ошибка при освобождении карта")
		}

		return nil
	})
}
