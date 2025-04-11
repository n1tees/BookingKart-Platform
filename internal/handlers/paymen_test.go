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

func setupPaymentRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	db.InitDB()
	r := gin.Default()

	r.GET("/payments", handlers.GetAllPayments)
	r.GET("/payments/:id", handlers.GetUserPayments)
	r.POST("/payments", handlers.CreatePayment)

	return r
}

func TestCreatePayment(t *testing.T) {
	r := setupPaymentRouter()

	body := map[string]interface{}{
		"user_id": 1,
		"amount":  500.0,
	}
	jsonValue, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "/payments", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "Баланс пополнен")
}

func TestGetAllPayments(t *testing.T) {
	r := setupPaymentRouter()
	req, _ := http.NewRequest("GET", "/payments", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetUserPayments(t *testing.T) {
	r := setupPaymentRouter()
	req, _ := http.NewRequest("GET", "/payments/1", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
