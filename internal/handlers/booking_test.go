package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/n1tees/BookingKart-Platform/internal/db"
	"github.com/n1tees/BookingKart-Platform/internal/handlers"
)

func setupBookingRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	db.InitDB()
	r := gin.Default()

	r.GET("/bookings", handlers.GetAllBookings)
	r.GET("/bookings/:id", handlers.GetBookingByID)
	r.POST("/bookings", handlers.CreateBooking)
	r.PATCH("/bookings/:id", handlers.PatchBooking)
	r.DELETE("/bookings/:id", handlers.DeleteBooking)

	return r
}

func TestCreateBooking(t *testing.T) {
	r := setupBookingRouter()

	body := map[string]interface{}{
		"user_id":    1,
		"kart_id":    1,
		"track_id":   1,
		"date":       "2025-05-01T00:00:00Z",
		"time_start": "2025-05-01T14:00:00Z",
		"duration":   20,
		"status":     "Запланировано",
	}
	jsonValue, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "/bookings", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "Бронирование создано")
}

func TestGetAllBookings(t *testing.T) {
	r := setupBookingRouter()
	req, _ := http.NewRequest("GET", "/bookings", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetBookingNotFound(t *testing.T) {
	r := setupBookingRouter()
	req, _ := http.NewRequest("GET", "/bookings/99999", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestDeleteBookingNotFound(t *testing.T) {
	r := setupBookingRouter()
	req, _ := http.NewRequest("DELETE", "/bookings/99999", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
