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
	"github.com/n1tees/BookingKart-Platform/internal/models"
)

func setupKartodromRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	db.InitDB() // подключение к тестовой БД
	r := gin.Default()

	r.GET("/kartodroms", handlers.GetAllKartodroms)
	r.POST("/kartodroms", handlers.CreateKartodrom)
	r.PUT("/kartodroms/:id", handlers.PatchKartodrom)
	r.PATCH("/kartodroms/:id", handlers.PatchKartodrom)
	r.DELETE("/kartodroms/:id", handlers.DeleteKartodrom)

	return r
}

func TestCreateKartodrom(t *testing.T) {
	router := setupKartodromRouter()

	body := models.Kartodrom{
		Name:     "Тестовый Карт",
		Location: "ТестГород",
		Phone:    "89990001122",
		Email:    "test@kart.ru",
	}
	jsonBody, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "/kartodroms", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestGetKartodroms(t *testing.T) {
	router := setupKartodromRouter()

	req, _ := http.NewRequest("GET", "/kartodroms", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestPatchKartodrom(t *testing.T) {
	router := setupKartodromRouter()

	// допустим, у нас уже есть ID = 1
	update := map[string]string{"Name": "Обновлённый Карт"}
	jsonUpdate, _ := json.Marshal(update)

	req, _ := http.NewRequest("PATCH", "/kartodroms/1", bytes.NewBuffer(jsonUpdate))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeleteKartodrom(t *testing.T) {
	router := setupKartodromRouter()

	req, _ := http.NewRequest("DELETE", "/kartodroms/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
