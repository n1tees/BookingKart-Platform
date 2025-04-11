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

func setupUserRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	db.InitDB() // подключение к тестовой БД
	r := gin.Default()

	r.GET("/users", handlers.GetAllUsers)
	r.GET("/users/:id", handlers.GetUserByID)
	r.POST("/users", handlers.CreateUser)
	r.PATCH("/users/:id", handlers.PatchUser)
	r.DELETE("/users/:id", handlers.DeleteUser)

	return r
}

func TestCreateUser(t *testing.T) {
	r := setupUserRouter()

	body := map[string]interface{}{
		"fname":    "Test",
		"phone":    "89998887766",
		"birthday": "2000-01-01T00:00:00Z",
		"password": "testpass",
	}
	jsonValue, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "Пользователь создан")
}

func TestGetAllUsers(t *testing.T) {
	r := setupUserRouter()
	req, _ := http.NewRequest("GET", "/users", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeleteUser_NotFound(t *testing.T) {
	r := setupUserRouter()
	req, _ := http.NewRequest("DELETE", "/users/99999", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}
