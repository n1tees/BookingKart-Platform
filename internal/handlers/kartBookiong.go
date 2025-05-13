package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/n1tees/BookingKart-Platform/internal/services"
)

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

func ActivateKartBookingHandler(c *gin.Context) {
	bookingID, _ := strconv.ParseUint(c.Param("bookingId"), 10, 64)
	kartID, _ := strconv.ParseUint(c.Param("kartId"), 10, 64)
	if err := services.ActivateKartBooking(uint(bookingID), uint(kartID)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "бронирование карты активировано"})
}

func FinishKartBookingHandler(c *gin.Context) {
	bookingID, _ := strconv.ParseUint(c.Param("bookingId"), 10, 64)
	kartID, _ := strconv.ParseUint(c.Param("kartId"), 10, 64)
	if err := services.FinishKartBooking(uint(bookingID), uint(kartID)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "бронирование карты завершено"})
}

func CancelKartBookingHandler(c *gin.Context) {
	bookingID, _ := strconv.ParseUint(c.Param("bookingId"), 10, 64)
	kartID, _ := strconv.ParseUint(c.Param("kartId"), 10, 64)
	if err := services.CancelKartBooking(uint(bookingID), uint(kartID)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "бронирование карты отменено"})
}
