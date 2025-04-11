package handlers

import (
	"net/http"

	"github.com/n1tees/BookingKart-Platform/internal/db"
	"github.com/n1tees/BookingKart-Platform/internal/models"

	"github.com/gin-gonic/gin"
)

// Create
func CreateKartodrom(c *gin.Context) {
	var kartodrom models.Kartodrom

	// Привязка тела запроса к структуре
	if err := c.ShouldBindJSON(&kartodrom); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Создание в базе
	if err := db.DB.Create(&kartodrom).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка создания картодрома"})
		return
	}

	// Успешный ответ
	c.JSON(http.StatusCreated, gin.H{
		"message": "Картодром добавлен",
		"id":      kartodrom.ID,
	})
}

// Get All
func GetAllKartodroms(c *gin.Context) {
	var kartodroms []models.Kartodrom

	// Получаем все записи из таблицы
	if err := db.DB.Find(&kartodroms).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения данных из БД"})
		return
	}

	// Возвращаем список картодромов в формате JSON
	c.JSON(http.StatusOK, kartodroms)
}

// Get By ID
func GetKartodromByID(c *gin.Context) {
	id := c.Param("id")
	var kartodrom models.Kartodrom
	if err := db.DB.First(&kartodrom, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Картодром не найден"})
		return
	}
	c.JSON(http.StatusOK, kartodrom)
}

// Patch
func PatchKartodrom(c *gin.Context) {
	id := c.Param("id")

	// Проверка, существует ли картодром
	var kartodrom models.Kartodrom
	if err := db.DB.First(&kartodrom, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Картодром не найден"})
		return
	}

	// Парсим тело запроса как map[string]interface{}
	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат данных"})
		return
	}

	// Обновляем только указанные поля
	if err := db.DB.Model(&kartodrom).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обновления"})
		return
	}

	// Отдаём результат
	c.JSON(http.StatusOK, gin.H{
		"message": "Картодром обновлён (частично)",
		"updated": updates,
	})
}

// Delete
func DeleteKartodrom(c *gin.Context) {
	id := c.Param("id")

	// Удаление по ID
	if err := db.DB.Delete(&models.Kartodrom{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка удаления картодрома"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Картодром удалён",
		"id":      id,
	})
}
