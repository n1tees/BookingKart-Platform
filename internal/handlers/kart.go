package handlers

import (
	"net/http"

	"github.com/n1tees/BookingKart-Platform/internal/db"
	"github.com/n1tees/BookingKart-Platform/internal/models"

	"github.com/gin-gonic/gin"
)

// Получить все карты
func GetAllKarts(c *gin.Context) {
	var karts []models.Kart
	if err := db.DB.Preload("Kartodrom").Preload("KartModel").Find(&karts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении картов"})
		return
	}
	c.JSON(http.StatusOK, karts)
}

// Получить карт по ID
func GetKartByID(c *gin.Context) {
	id := c.Param("id")
	var kart models.Kart
	if err := db.DB.Preload("Kartodrom").Preload("KartModel").First(&kart, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Карт не найден"})
		return
	}
	c.JSON(http.StatusOK, kart)
}

// Создание карта
func CreateKart(c *gin.Context) {
	type Input struct {
		KartodromID uint              `json:"kartodrom_id"`
		KartModelID uint              `json:"kart_model_id"`
		KartStatus  models.KartStatus `json:"kart_status"`
	}

	var input Input
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	kart := models.Kart{
		KartodromID: input.KartodromID,
		KartModelID: input.KartModelID,
		KartStatus:  input.KartStatus,
	}

	if err := db.DB.Create(&kart).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при создании карта"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Карт создан", "id": kart.ID})
}

// Частичное обновление карта
func PatchKart(c *gin.Context) {
	id := c.Param("id")
	var kart models.Kart

	if err := db.DB.First(&kart, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Карт не найден"})
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.DB.Model(&kart).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при обновлении карта"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Карт обновлён", "updated": updates})
}

// Удаление карта
func DeleteKart(c *gin.Context) {
	id := c.Param("id")
	if err := db.DB.Delete(&models.Kart{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при удалении карта"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Карт удалён"})
}
