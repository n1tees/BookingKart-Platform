package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/n1tees/BookingKart-Platform/pkg/models"
	"github.com/n1tees/BookingKart-Platform/pkg/services"
)

type BookingRequest struct {
	TrackID     uint               `json:"track_id"`
	CustomerID  uint               `json:"customer_id"`
	Date        string             `json:"date"`
	StartTime   string             `json:"start_time"`
	Duration    uint               `json:"duration"`
	BookingType models.BookingType `json:"booking_type"`
	RiderCount  uint               `json:"rider_count"`
}

// CreateBookingHandler godoc
// @Summary Создать бронирование
// @Tags booking
// @Accept json
// @Produce json
// @Param input body handlers.BookingRequest true "Данные бронирования"
// @Success 201 {object} map[string]uint
// @Failure 400 {object} map[string]string
// @Router /bookings [post]
func CreateBookingHandler(c *gin.Context) {
	var req BookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный формат запроса"})
		return
	}

	parsedDate, err := MakeDateByString(req.Date)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный формат даты (YYYY-MM-DD)"})
		return
	}
	parsedTime, err := MakeTimeByString(req.StartTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный формат времени (HH:MM)"})
		return
	}

	input := services.BookingInput{
		TrackID:     req.TrackID,
		CustomerID:  req.CustomerID,
		Date:        models.LocalTime{Time: parsedDate},
		StartTime:   models.LocalTime{Time: parsedTime},
		Duration:    req.Duration,
		BookingType: req.BookingType,
		RiderCount:  req.RiderCount,
	}

	bookingID, err := services.ReserveBooking(input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"booking_id": bookingID})
}

// ActivateBookingHandler godoc
// @Summary Активировать бронирование
// @Tags booking
// @Param id path int true "ID бронирования"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /bookings/{id}/activate [post]
func ActivateBookingHandler(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := services.ActivateBooking(uint(id)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "бронирование активировано"})
}

// CloseBookingHandler godoc
// @Summary Завершить бронирование
// @Tags booking
// @Param id path int true "ID бронирования"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /bookings/{id}/close [post]
func CloseBookingHandler(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := services.CloseBooking(uint(id)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "бронирование завершено"})
}

// CancelBookingHandler godoc
// @Summary Отменить бронирование
// @Tags booking
// @Param id path int true "ID бронирования"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /bookings/{id}/cancel [post]
func CancelBookingHandler(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := services.CancelBooking(uint(id)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "бронирование отменено"})
}

// GetBookingsByDateHandler godoc
// @Summary Получить список бронирований по дате
// @Tags booking
// @Param id path int true "ID картодрома"
// @Param date query string true "Дата (YYYY-MM-DD)"
// @Success 200 {array} models.Booking
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /kartodrom/{id}/bookings [get]

func GetBookingsByDateHandler(c *gin.Context) {
	kartodromID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	dateStr := c.Query("date")
	parsedDate, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный формат даты (YYYY-MM-DD)"})
		return
	}
	bookings, err := services.GetBookingsByDate(uint(kartodromID), parsedDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, bookings)
}
