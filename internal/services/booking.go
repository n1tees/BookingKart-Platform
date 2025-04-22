package services

import (
	"errors"
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
	StartTime   time.Time
	Duration    uint // в минутах
	BookingType models.BookingType
	RiderCount  uint // если пусто → ставим 1
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
		if riderCount == 0 {
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
			TrackID:     input.TrackID,
			CustomerID:  input.CustomerID,
			Date:        date,
			StartTime:   input.StartTime,
			EndTime:     input.StartTime.Add(time.Duration(input.Duration) * time.Minute),
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
			return errors.New("бронирование не найдено")
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

		if booking.Status != models.BookingActive {
			return errors.New("бронирование уже не активно и не может быть отменено")
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
	// Получаем ID всех треков в картодроме
	var trackIDs []uint
	if err := db.DB.Model(&models.Track{}).
		Where("kartodrom_id = ?", kartodromID).
		Pluck("id", &trackIDs).Error; err != nil {
		return nil, errors.New("ошибка при получении треков картодрома")
	}

	if len(trackIDs) == 0 {
		return &[]models.Booking{}, nil
	}

	// Только дата
	dateOnly := time.Date(targetDate.Year(), targetDate.Month(), targetDate.Day(), 0, 0, 0, 0, targetDate.Location())

	// Получаем бронирования на дату с нужными статусами
	var bookings []models.Booking
	if err := db.DB.
		Where("date = ? AND status IN ? AND track_id IN ?", dateOnly, []models.BookingStatus{models.BookingActive, models.BookingReserve}, trackIDs).
		Find(&bookings).Error; err != nil {
		return nil, errors.New("ошибка при получении бронирований")
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
