package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/n1tees/BookingKart-Platform/internal/db"
	"github.com/n1tees/BookingKart-Platform/internal/models"
)

// Получить список всех треков
func GetAllTracks(c *gin.Context) {
	var tracks []models.Track
	if err := db.DB.Preload("Kartodrom").Find(&tracks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении треков"})
		return
	}
	c.JSON(http.StatusOK, tracks)
}

// Получить трек по ID
func GetTrackByID(c *gin.Context) {
	id := c.Param("id")
	var track models.Track
	if err := db.DB.Preload("Kartodrom").First(&track, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Трек не найден"})
		return
	}
	c.JSON(http.StatusOK, track)
}

// Создание трека
func CreateTrack(c *gin.Context) {
	type Input struct {
		KartodromID    uint            `json:"kartodrom_id"`
		Name           string          `json:"name"`
		Length         uint            `json:"length"`
		DifLevel       models.DifLevel `json:"dif_level"`
		PricePerMinute float64         `json:"price_per_minute"`
		MaxKart        uint            `json:"max_kart"`
	}

	var input Input
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	track := models.Track{
		KartodromID:    input.KartodromID,
		Name:           input.Name,
		Length:         input.Length,
		DifLevel:       input.DifLevel,
		PricePerMinute: input.PricePerMinute,
		MaxKart:        input.MaxKart,
	}

	if err := db.DB.Create(&track).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при создании трека"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Трек создан", "id": track.ID})
}

// Обновление трека
func PatchTrack(c *gin.Context) {
	id := c.Param("id")
	var track models.Track

	if err := db.DB.First(&track, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Трек не найден"})
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.DB.Model(&track).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при обновлении трека"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Трек обновлён", "updated": updates})
}

// Удаление трека
func DeleteTrack(c *gin.Context) {
	id := c.Param("id")
	if err := db.DB.Delete(&models.Track{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при удалении трека"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Трек удалён"})
}
