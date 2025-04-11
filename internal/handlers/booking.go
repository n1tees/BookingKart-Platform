package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/n1tees/BookingKart-Platform/internal/db"
	"github.com/n1tees/BookingKart-Platform/internal/models"
)

// Получить все бронирования
func GetAllBookings(c *gin.Context) {
	var bookings []models.Booking
	if err := db.DB.Preload("User").Preload("Kart").Preload("Track").Find(&bookings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении бронирований"})
		return
	}
	c.JSON(http.StatusOK, bookings)
}

// Получить бронирование по ID
func GetBookingByID(c *gin.Context) {
	id := c.Param("id")
	var booking models.Booking
	if err := db.DB.Preload("User").Preload("Kart").Preload("Track").First(&booking, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Бронирование не найдено"})
		return
	}
	c.JSON(http.StatusOK, booking)
}

// Создание бронирования
func CreateBooking(c *gin.Context) {
	type Input struct {
		UserID      uint                 `json:"user_id"`
		TrackID     uint                 `json:"track_id"`
		CustomerID  uint                 `json:"customer_id"`
		Date        time.Time            `json:"date"`
		TimeStart   time.Time            `json:"time_start"`
		Duration    int                  `json:"duration"` // в минутах
		BookingType models.BookingType   `json:"booking_type"`
		Status      models.BookingStatus `json:"status"`
	}

	var input Input
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Получаем цену с трека
	var track models.Track
	if err := db.DB.First(&track, input.TrackID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Трек не найден"})
		return
	}

	total := float64(input.Duration) * track.PricePerMinute
	timeEnd := input.TimeStart.Add(time.Duration(input.Duration) * time.Minute)

	booking := models.Booking{
		TrackID:     input.TrackID,
		CustomerID:  input.CustomerID,
		Date:        input.Date,
		TimeStart:   input.TimeStart,
		TimeEnd:     timeEnd,
		Duration:    uint(input.Duration),
		BookingType: input.BookingType,
		Status:      input.Status,
		TotalPrice:  total,
	}

	if err := db.DB.Create(&booking).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка создания бронирования"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Бронирование создано", "id": booking.ID})
}

// Обновление бронирования
func PatchBooking(c *gin.Context) {
	id := c.Param("id")
	var booking models.Booking

	if err := db.DB.First(&booking, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Бронирование не найдено"})
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.DB.Model(&booking).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при обновлении бронирования"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Бронирование обновлено", "updated": updates})
}

// ✅ Удаление бронирования
func DeleteBooking(c *gin.Context) {
	id := c.Param("id")
	if err := db.DB.Delete(&models.Booking{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при удалении бронирования"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Бронирование удалено"})
}
