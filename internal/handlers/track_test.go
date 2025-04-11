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

func setupTrackRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	db.InitDB()
	r := gin.Default()

	r.GET("/tracks", handlers.GetAllTracks)
	r.GET("/tracks/:id", handlers.GetTrackByID)
	r.POST("/tracks", handlers.CreateTrack)
	r.PATCH("/tracks/:id", handlers.PatchTrack)
	r.DELETE("/tracks/:id", handlers.DeleteTrack)

	return r
}

func TestCreateTrack(t *testing.T) {
	r := setupTrackRouter()

	body := map[string]interface{}{
		"kartodrom_id":     1,
		"name":             "Тестовый трек",
		"length":           800,
		"dif_level":        "Средний",
		"price_per_minute": 30.5,
		"max_kart":         10,
	}
	jsonValue, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "/tracks", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "Трек создан")
}

func TestGetAllTracks(t *testing.T) {
	r := setupTrackRouter()
	req, _ := http.NewRequest("GET", "/tracks", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetTrackNotFound(t *testing.T) {
	r := setupTrackRouter()
	req, _ := http.NewRequest("GET", "/tracks/99999", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestDeleteTrackNotFound(t *testing.T) {
	r := setupTrackRouter()
	req, _ := http.NewRequest("DELETE", "/tracks/99999", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
