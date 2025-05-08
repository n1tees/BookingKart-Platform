package handlers

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/n1tees/BookingKart-Platform/internal/services"
)

func RegisterHandler(c *gin.Context) {
	var input services.RegisterInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный формат запроса"})
		return
	}

	// Парсим дату рождения
	birth, err := MakeDateByString(input.BirthDay)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный формат даты. Ожидается YYYY-MM-DD"})
		return
	}

	// Передаём дату в структуру
	input.ParsedBirthDay = birth

	userID, err := services.RegUser(input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"user_id": userID})
}

func LoginHandler(c *gin.Context) {
	var input services.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный формат запроса"})
		return
	}

	token, err := services.LoginUser(input)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

// функции для работы с датами и временем
func MakeDateByString(date string) (time.Time, error) {
	parsedDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		return time.Time{}, errors.New("ошибка при парсинге даты")
	}
	return parsedDate, nil
}

func MakeTimeByString(timeStr string) (time.Time, error) {
	parsedTime, err := time.Parse("15:04", timeStr)
	if err != nil {
		return time.Time{}, errors.New("ошибка при парсинге времени")
	}
	return parsedTime, nil
}
