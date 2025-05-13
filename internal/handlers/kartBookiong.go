package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/n1tees/BookingKart-Platform/internal/services"
)

// GetAvailableKartsForBookingHandler godoc
// @Summary Получить доступные карты для бронирования
// @Tags kartbooking
// @Param id path int true "ID картодрома"
// @Param start query string true "Дата начала (2025-05-15T10:00:00)"
// @Param end query string true "Дата окончания (2025-05-15T11:00:00)"
// @Success 200 {array} models.Kart
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /kartodrom/{id}/karts [get]
func GetAvailableKartsForBookingHandler(c *gin.Context) {
	kartodromID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	startStr := c.Query("start")
	endStr := c.Query("end")
	start, err1 := time.Parse(time.RFC3339, startStr)
	end, err2 := time.Parse(time.RFC3339, endStr)
	if err1 != nil || err2 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный формат даты. Пример: 2025-05-15T10:00:00"})
		return
	}

	karts, err := services.GetAvailableKartsForBooking(uint(kartodromID), start, end)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, karts)
}

type KartBookingRequest struct {
	BookingID uint `json:"booking_id"`
	KartID    uint `json:"kart_id"`
}

// ReserveKartHandler godoc
// @Summary Зарезервировать карт для бронирования
// @Tags kartbooking
// @Param booking body handlers.KartBookingRequest true "Данные бронирования карта"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /kartbookings [post]
func ReserveKartHandler(c *gin.Context) {
	var req KartBookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный формат запроса"})
		return
	}

	if err := services.ReserveKart(req.BookingID, req.KartID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "карт зарезервирован"})
}

// ActivateKartBookingHandler godoc
// @Summary Активировать бронирование карта
// @Tags kartbooking
// @Param bookingId path int true "ID бронирования"
// @Param kartId path int true "ID карта"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /kartbookings/{bookingId}/{kartId}/activate [post]
func ActivateKartBookingHandler(c *gin.Context) {
	bookingID, _ := strconv.ParseUint(c.Param("bookingId"), 10, 64)
	kartID, _ := strconv.ParseUint(c.Param("kartId"), 10, 64)
	if err := services.ActivateKartBooking(uint(bookingID), uint(kartID)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "бронирование карты активировано"})
}

// FinishKartBookingHandler godoc
// @Summary Завершить бронирование карта
// @Tags kartbooking
// @Param bookingId path int true "ID бронирования"
// @Param kartId path int true "ID карта"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /kartbookings/{bookingId}/{kartId}/finish [post]
func FinishKartBookingHandler(c *gin.Context) {
	bookingID, _ := strconv.ParseUint(c.Param("bookingId"), 10, 64)
	kartID, _ := strconv.ParseUint(c.Param("kartId"), 10, 64)
	if err := services.FinishKartBooking(uint(bookingID), uint(kartID)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "бронирование карты завершено"})
}

// CancelKartBookingHandler godoc
// @Summary Отменить бронирование карта
// @Tags kartbooking
// @Param bookingId path int true "ID бронирования"
// @Param kartId path int true "ID карта"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /kartbookings/{bookingId}/{kartId}/cancel [post]
func CancelKartBookingHandler(c *gin.Context) {
	bookingID, _ := strconv.ParseUint(c.Param("bookingId"), 10, 64)
	kartID, _ := strconv.ParseUint(c.Param("kartId"), 10, 64)
	if err := services.CancelKartBooking(uint(bookingID), uint(kartID)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "бронирование карты отменено"})
}
