package services

import (
	"github.com/n1tees/BookingKart-Platform/pkg/db"
	"github.com/n1tees/BookingKart-Platform/pkg/models"
	"gorm.io/gorm"

	"errors"
	"time"
)

type AvailableTrack struct {
	ID          uint
	Name        string
	Length      uint
	DifLevel    models.DifLevel
	PricePerMin float64
	MaxKarts    uint
	FreeSlots   uint
}

// Получить доступные треки на картодроме
func GetAvailableTracks(kartodromID uint) (*[]AvailableTrack, error) {
	var allTracks []models.Track
	if err := db.DB.Where("kartodrom_id = ?", kartodromID).Find(&allTracks).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("треки не найдены")
		}
		return nil, errors.New("ошибка при поиске треков")
	}

	var result []AvailableTrack
	for _, track := range allTracks {
		count, err := сountRidersOnTrack(track.ID)
		if err != nil {
			return nil, err
		}

		if count < track.MaxKarts {
			result = append(result, AvailableTrack{
				ID:          track.ID,
				Name:        track.Name,
				Length:      track.Length,
				DifLevel:    track.DifLevel,
				PricePerMin: track.PricePerMin,
				MaxKarts:    track.MaxKarts,
				FreeSlots:   track.MaxKarts - count,
			})
		}
	}

	return &result, nil
}

// Получить трек по id
func GetTrackByID(trackID uint) (*models.Track, error) {

	var track models.Track

	if err := db.DB.First(&track, trackID).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("трек не найден")

		} else {
			return nil, errors.New("ошибка при поиске трека")
		}
	}

	return &track, nil
}

// посчитать количество райдеров на треке
func сountRidersOnTrack(trackID uint) (uint, error) {
	now := time.Now()

	var totalRiders uint

	err := db.DB.Model(&models.Booking{}).
		Select("COALESCE(SUM(rider_count), 0)").
		Where("track_id = ? AND status = ? AND start_time <= ? AND end_time >= ?",
			trackID, models.BookingActive, now, now).
		Scan(&totalRiders).Error

	if err != nil {
		return 0, errors.New("ошибка при подсчёте райдеров на треке")
	}

	return totalRiders, nil
}
