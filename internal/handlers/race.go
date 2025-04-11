package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/n1tees/BookingKart-Platform/internal/db"
	"github.com/n1tees/BookingKart-Platform/internal/models"
)

// Получить все гонки
func GetAllRaces(c *gin.Context) {
	var races []models.Race
	if err := db.DB.Preload("Track").Find(&races).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении гонок"})
		return
	}
	c.JSON(http.StatusOK, races)
}

// Получить гонку по ID
func GetRaceByID(c *gin.Context) {
	id := c.Param("id")
	var race models.Race
	if err := db.DB.Preload("Track").First(&race, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Гонка не найдена"})
		return
	}
	c.JSON(http.StatusOK, race)
}

// Создание гонки
func CreateRace(c *gin.Context) {
	type Input struct {
		TrackID    uint              `json:"track_id"`
		Date       time.Time         `json:"date"`
		TimeStart  time.Time         `json:"time_start"`
		Duration   uint              `json:"duration"`
		Laps       uint              `json:"laps"`
		RaceType   models.RaceType   `json:"race_type"`
		RaceStatus models.RaceStatus `json:"race_status"`
		TotalPrice float64           `json:"total_price"`
	}

	var input Input
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	timeEnd := input.TimeStart.Add(time.Duration(input.Duration) * time.Minute)

	race := models.Race{
		TrackID:    input.TrackID,
		Date:       input.Date,
		TimeStart:  input.TimeStart,
		TimeEnd:    timeEnd,
		Duration:   input.Duration,
		Laps:       input.Laps,
		RaceType:   input.RaceType,
		RaceStatus: input.RaceStatus,
		TotalPrice: input.TotalPrice,
	}

	if err := db.DB.Create(&race).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при создании гонки"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Гонка создана", "id": race.ID})
}

// Обновление гонки
func PatchRace(c *gin.Context) {
	id := c.Param("id")
	var race models.Race

	if err := db.DB.First(&race, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Гонка не найдена"})
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.DB.Model(&race).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при обновлении гонки"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Гонка обновлена", "updated": updates})
}

// Удаление гонки
func DeleteRace(c *gin.Context) {
	id := c.Param("id")
	if err := db.DB.Delete(&models.Race{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при удалении гонки"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Гонка удалена"})
}
