package services

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/n1tees/BookingKart-Platform/internal/db"
	"github.com/n1tees/BookingKart-Platform/internal/models"
	"gorm.io/gorm"
)

// Стуктура ввода для бронирования
type BookingInput struct {
	TrackID     uint
	CustomerID  uint
	Date        models.LocalTime
	StartTime   models.LocalTime
	Duration    uint
	BookingType models.BookingType
	RiderCount  uint
}

// Создание бронирования
func ReserveBooking(input BookingInput) (uint, error) {

	var mu sync.Mutex
	mu.Lock()
	defer mu.Unlock()

	var newBooking models.Booking

	err := db.DB.Transaction(func(tx *gorm.DB) error {

		var track models.Track
		if err := tx.First(&track, input.TrackID).Error; err != nil {
			return errors.New("трек не найден")
		}

		riderCount := input.RiderCount
		if riderCount <= 0 {
			riderCount = 1
		}

		totalPrice := float64(input.Duration) * track.PricePerMin * float64(riderCount)

		err := ChargeFromBalance(input.CustomerID, totalPrice)
		if err != nil {
			return err
		}

		date := time.Date(
			input.StartTime.Year(), input.StartTime.Month(), input.StartTime.Day(),
			0, 0, 0, 0, input.StartTime.Location(),
		)

		newBooking = models.Booking{
			TrackID:    input.TrackID,
			CustomerID: input.CustomerID,
			Date:       date,
			StartTime:  input.StartTime,
			EndTime: models.LocalTime{
				Time: input.StartTime.Time.Add(time.Duration(input.Duration) * time.Minute),
			},
			Duration:    input.Duration,
			TotalPrice:  totalPrice,
			BookingType: input.BookingType,
			Status:      models.BookingReserve,
			RiderCount:  riderCount,
		}

		if err := tx.Create(&newBooking).Error; err != nil {
			return errors.New("ошибка при создании бронирования")
		}

		return nil
	})

	if err != nil {
		return 0, err
	}

	return newBooking.ID, nil
}

// Активация бронирования
func ActivateBooking(bookingID uint) error {
	var mu sync.Mutex
	mu.Lock()
	defer mu.Unlock()

	return db.DB.Transaction(func(tx *gorm.DB) error {
		var booking models.Booking
		if err := tx.First(&booking, bookingID).Error; err != nil {

			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("бронирование не найдено")

			} else {
				fmt.Print(err.Error())
				return errors.New("ошибка при поиске бронирования")
			}
		}

		if booking.Status != models.BookingReserve {
			return errors.New("только забронированное бронирование может быть активировано")
		}

		if err := tx.Model(&booking).Update("status", models.BookingActive).Error; err != nil {
			return errors.New("ошибка при активации бронирования")
		}

		return nil
	})
}

// Завершение бронирования
func CloseBooking(bookingID uint) error {
	var mu sync.Mutex
	mu.Lock()
	defer mu.Unlock()

	return db.DB.Transaction(func(tx *gorm.DB) error {
		var booking models.Booking
		if err := tx.First(&booking, bookingID).Error; err != nil {
			return errors.New("бронирование не найдено")
		}

		if booking.Status != models.BookingActive {
			return errors.New("только активное бронирование может быть завершено")
		}

		if err := tx.Model(&booking).Update("status", models.BookingClose).Error; err != nil {
			return errors.New("ошибка при завершении бронирования")
		}

		return nil
	})
}

// Отменить бронь
func CancelBooking(bookingID uint) error {
	var mu sync.Mutex
	mu.Lock()
	defer mu.Unlock()

	err := db.DB.Transaction(func(tx *gorm.DB) error {
		var booking models.Booking
		if err := tx.First(&booking, bookingID).Error; err != nil {
			return errors.New("бронирование не найдено")
		}

		if booking.Status == models.BookingActive || booking.Status == models.BookingClose {
			return errors.New("бронирование не может быть отменено")
		}

		// Меняем статус на Cancelled
		if err := tx.Model(&booking).Update("status", models.BookingCancel).Error; err != nil {
			return errors.New("ошибка при отмене бронирования")
		}

		// Возвращаем деньги
		if err := RefundToBalance(booking.CustomerID, booking.TotalPrice); err != nil {
			return errors.New("ошибка при возврате средств")
		}

		return nil
	})

	return err
}

// Получить бронирования на текущий день
func GetBookingsByDate(kartodromID uint, targetDate time.Time) (*[]models.Booking, error) {

	var trackIDs []uint
	if err := db.DB.Model(&models.Track{}).
		Where("kartodrom_id = ?", kartodromID).
		Pluck("id", &trackIDs).Error; err != nil {
		return nil, fmt.Errorf("ошибка при получении треков картодрома: %w", err)
	}

	if len(trackIDs) == 0 {
		return nil, errors.New("нет треков")
	}

	// Форматирование даты в YYYY-MM-DD
	dateStr := targetDate.Format("2006-01-02")

	var bookings []models.Booking
	if err := db.DB.
		Where("DATE(date) = ? AND status IN ? AND track_id IN ?", dateStr,
			[]models.BookingStatus{models.BookingActive, models.BookingReserve}, trackIDs).
		Find(&bookings).Error; err != nil {
		return nil, fmt.Errorf("ошибка при получении бронирований: %w", err)
	}

	return &bookings, nil
}

// Функционал:
// 1 Получить расписание броней на текущий день для определенного картодрома
// 2 Работа со статусами
// 	 - Создать бронь
// 	 - Активировать бронь
// 	 - Закрыть бронь
// 	 - Отменить бронь
