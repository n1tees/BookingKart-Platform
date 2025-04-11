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

func setupRaceRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	db.InitDB()
	r := gin.Default()

	r.GET("/races", handlers.GetAllRaces)
	r.GET("/races/:id", handlers.GetRaceByID)
	r.POST("/races", handlers.CreateRace)
	r.PATCH("/races/:id", handlers.PatchRace)
	r.DELETE("/races/:id", handlers.DeleteRace)

	return r
}

func TestCreateRace(t *testing.T) {
	r := setupRaceRouter()

	body := map[string]interface{}{
		"track_id":    1,
		"date":        "2025-06-01T00:00:00Z",
		"time_start":  "2025-06-01T12:00:00Z",
		"duration":    30,
		"laps":        10,
		"race_type":   "спринт",
		"race_status": "Ожидается",
		"total_price": 1500.0,
	}
	jsonValue, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "/races", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "Гонка создана")
}

func TestGetAllRaces(t *testing.T) {
	r := setupRaceRouter()
	req, _ := http.NewRequest("GET", "/races", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetRaceNotFound(t *testing.T) {
	r := setupRaceRouter()
	req, _ := http.NewRequest("GET", "/races/99999", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestDeleteRaceNotFound(t *testing.T) {
	r := setupRaceRouter()
	req, _ := http.NewRequest("DELETE", "/races/99999", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
