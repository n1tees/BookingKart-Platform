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

func setupKartRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	db.InitDB()
	r := gin.Default()

	r.GET("/karts", handlers.GetAllKarts)
	r.GET("/karts/:id", handlers.GetKartByID)
	r.POST("/karts", handlers.CreateKart)
	r.PATCH("/karts/:id", handlers.PatchKart)
	r.DELETE("/karts/:id", handlers.DeleteKart)

	return r
}

func TestCreateKart(t *testing.T) {
	r := setupKartRouter()

	body := map[string]interface{}{
		"kartodrom_id":  1,
		"kart_model_id": 1,
		"kart_status":   "Доступен",
	}
	jsonValue, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "/karts", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "Карт создан")
}

func TestGetAllKarts(t *testing.T) {
	r := setupKartRouter()
	req, _ := http.NewRequest("GET", "/karts", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetKartNotFound(t *testing.T) {
	r := setupKartRouter()
	req, _ := http.NewRequest("GET", "/karts/99999", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestDeleteKartNotFound(t *testing.T) {
	r := setupKartRouter()
	req, _ := http.NewRequest("DELETE", "/karts/99999", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
