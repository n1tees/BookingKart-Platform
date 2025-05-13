package services

import (
	"errors"
	"sync"
	"time"

	"github.com/n1tees/BookingKart-Platform/internal/db"
	"github.com/n1tees/BookingKart-Platform/internal/models"
	"gorm.io/gorm"
)

// Стуктура ввода для заезда
type RaceInput struct {
	TrackID   uint
	Date      string
	TimeStart string
	RaceType  models.RaceType
	Laps      uint
	Duration  uint
	Status    models.RaceStatus

	ParsedDate      time.Time `json:"-"`
	ParsedTimeStart time.Time `json:"-"`
}

// Создание заезда
func CreateRace(input RaceInput) (uint, error) {
	var mu sync.Mutex
	mu.Lock()
	defer mu.Unlock()

	var newRace models.Race

	err := db.DB.Transaction(func(tx *gorm.DB) error {
		// Проверяем трек
		var track models.Track
		if err := tx.First(&track, input.TrackID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("трек не найден")

			} else {
				return errors.New("ошибка при поиске трека")
			}
		}

		newRace = models.Race{
			TrackID:    input.TrackID,
			Date:       input.ParsedDate,
			TimeStart:  models.LocalTime{Time: input.ParsedTimeStart},
			Laps:       input.Laps,
			Duration:   input.Duration,
			TotalPrice: track.PricePerMin * float64(input.Duration),
			RaceType:   input.RaceType,
			Status:     models.RaceCreate,
		}

		if err := tx.Create(&newRace).Error; err != nil {
			return errors.New("ошибка при создании заезда")
		}

		return nil
	})

	if err != nil {
		return 0, err
	}

	return newRace.ID, nil
}

// Старт заезда
func StartRace(raceID uint) error {
	var mu sync.Mutex
	mu.Lock()
	defer mu.Unlock()

	return db.DB.Transaction(func(tx *gorm.DB) error {
		var race models.Race
		if err := tx.First(&race, raceID).Error; err != nil {

			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("заезд не найден")

			} else {
				return errors.New("ошибка при поиске заезда")
			}
		}

		if race.Status != models.RaceCreate {
			return errors.New("можно стартовать только созданную гонку")
		}

		if err := tx.Model(&race).Update("status", models.RaceStart).Error; err != nil {
			return errors.New("ошибка при старте гонки")
		}

		return nil
	})
}

// Завершение заезда
func FinishRace(raceID uint) error {
	var mu sync.Mutex
	mu.Lock()
	defer mu.Unlock()

	return db.DB.Transaction(func(tx *gorm.DB) error {
		var race models.Race
		if err := tx.First(&race, raceID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("заезд не найден")

			} else {
				return errors.New("ошибка при поиске заезда")
			}
		}

		if race.Status != models.RaceStart {
			return errors.New("можно завершить только начатую гонку")
		}

		if err := tx.Model(&race).Update("status", models.RaceFinish).Error; err != nil {
			return errors.New("ошибка при завершении гонки")
		}

		return nil
	})
}

// Отмена заезда
func CancelRace(raceID uint) error {
	var mu sync.Mutex
	mu.Lock()
	defer mu.Unlock()

	return db.DB.Transaction(func(tx *gorm.DB) error {
		var race models.Race
		if err := tx.First(&race, raceID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("заезд не найден")

			} else {

				return errors.New("ошибка при поиске заезда")
			}
		}

		if race.Status != models.RaceCreate {
			return errors.New("можно отменить только созданную гонку")
		}

		if err := tx.Model(&race).Update("status", models.RaceCanceled).Error; err != nil {
			return errors.New("ошибка при отмене гонки")
		}

		return nil
	})
}

// universeal
func UpdateRaceStatus(raceID uint, newStatus models.RaceStatus) error {
	var mu sync.Mutex
	mu.Lock()
	defer mu.Unlock()

	return db.DB.Transaction(func(tx *gorm.DB) error {
		var race models.Race
		if err := tx.First(&race, raceID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("заезд не найден")

			} else {
				return errors.New("ошибка при поиске заезда")
			}
		}

		switch newStatus {
		case models.RaceStart:
			if race.Status != models.RaceCreate {
				return errors.New("только созданную гонку можно стартовать")
			}
		case models.RaceFinish:
			if race.Status != models.RaceStart {
				return errors.New("только начатую гонку можно завершить")
			}
		case models.RaceCanceled:
			if race.Status != models.RaceCreate {
				return errors.New("только созданную гонку можно отменить")
			}
		default:
			return errors.New("неподдерживаемый целевой статус гонки")
		}

		if err := tx.Model(&race).Update("status", newStatus).Error; err != nil {
			return errors.New("ошибка при обновлении статуса гонки")
		}

		return nil
	})
}
