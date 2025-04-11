package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/n1tees/BookingKart-Platform/internal/db"
	"github.com/n1tees/BookingKart-Platform/internal/models"
	"gorm.io/gorm"
)

// Получить список всех пополнений
func GetAllPayments(c *gin.Context) {
	var payments []models.Payment
	if err := db.DB.Find(&payments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении платежей"})
		return
	}
	c.JSON(http.StatusOK, payments)
}

// Получить пополнения конкретного пользователя
func GetUserPayments(c *gin.Context) {
	id := c.Param("id")
	var payments []models.Payment
	if err := db.DB.Where("user_id = ?", id).Find(&payments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении платежей пользователя"})
		return
	}
	c.JSON(http.StatusOK, payments)
}

// Создание пополнения (и обновление баланса)
func CreatePayment(c *gin.Context) {
	type Input struct {
		UserID uint    `json:"user_id"`
		Amount float64 `json:"amount"`
	}

	var input Input
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	payment := models.Payment{
		UserID: input.UserID,
		Amount: input.Amount,
		Date:   time.Now(),
	}

	// Добавляем запись о пополнении
	if err := db.DB.Create(&payment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при создании пополнения"})
		return
	}

	// Обновляем баланс пользователя
	if err := db.DB.Model(&models.UserInfo{}).Where("id = ?", input.UserID).
		UpdateColumn("balance", gorm.Expr("balance + ?", input.Amount)).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Платёж записан, но ошибка при обновлении баланса"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Баланс пополнен", "payment_id": payment.ID})
}
