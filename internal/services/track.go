package services

import (
	"github.com/n1tees/BookingKart-Platform/internal/db"
	"github.com/n1tees/BookingKart-Platform/internal/models"

	"errors"
	"time"
)

// Получить доступные треки на картодроме
func GetAvailableTracks(kartodromID uint) ([]models.Track, error) {
	var allTracks []models.Track

	// Сначала получаем все треки картодрома
	if err := db.DB.Where("kartodrom_id = ?", kartodromID).Find(&allTracks).Error; err != nil {
		return nil, errors.New("ошибка при получении треков картодрома")
	}

	var availableTracks []models.Track

	// Фильтрация по количеству райдеров
	for _, track := range allTracks {
		count, err := CountRidersOnTrack(track.ID)
		if err != nil {
			return nil, err
		}

		if count < track.MaxKarts {
			availableTracks = append(availableTracks, track)
		}
	}

	return availableTracks, nil
}

// Посчитать количество райдеров на треке
func CountRidersOnTrack(trackID uint) (uint, error) {

	now := time.Now()

	var count int64

	if err := db.DB.Model(&models.Booking{}).
		Where("track_id = ? AND status = ? AND start_time <= ? AND end_time >= ?",
			trackID, models.BookingActive, now, now).
		Count(&count).Error; err != nil {
		return 0, errors.New("ошибка при подсчёте активных райдеров на треке")
	}

	return uint(count), nil
}

// Получить трек по id
func GetTrackByID(trackID uint) (*models.Track, error) {
	var track models.Track

	if err := db.DB.First(&track, trackID).Error; err != nil {
		return nil, errors.New("трек не найден")
	}

	return &track, nil
}
