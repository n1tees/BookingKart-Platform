package auth

import (
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"github.com/n1tees/BookingKart-Platform/config"
	"github.com/n1tees/BookingKart-Platform/internal/db"
	"github.com/n1tees/BookingKart-Platform/internal/models"
)

// структура запроса
type LoginInput struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// POST /auth/login
func LoginHandler(c *gin.Context) {

	//проверяем, корректен ли полученный ввод
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат входных данных"})
		return
	}

	// ищем логин в бд
	var auth models.AuthCredential
	if err := db.DB.Where("login = ?", input.Login).First(&auth).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный логин или пароль"})
		return
	}

	// сравниваем введённый пароль с захэшированным паролем из БД
	if err := bcrypt.CompareHashAndPassword(auth.PasswordHash, []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный логин или пароль"})
		return
	}

	// получаем тип пользователя через таблицу User
	var user models.User
	if err := db.DB.Where("auth_id = ?", auth.ID).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Пользователь не найден"})
		return
	}

	// генерация JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":   user.ID,
		"user_type": user.UserType,
		"exp":       time.Now().Add(time.Hour * 2).Unix(),
	})

	tokenString, err := token.SignedString([]byte(config.GetJWTSecret()))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при генерации токена"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}
