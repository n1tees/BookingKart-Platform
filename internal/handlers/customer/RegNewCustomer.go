package customer

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"github.com/n1tees/BookingKart-Platform/internal/db"
	"github.com/n1tees/BookingKart-Platform/internal/models"
)

// структура запроса
// минимальная информация для регистрации

type RegNewCustomerInput struct {
	Phone    string    `json:"phone" binding:"required"`
	FName    string    `json:"fname" binding:"required"`
	BirthDay time.Time `json:"birth_day" binding:"required"`
	Login    string    `json:"login" binding:"required"`
	Password string    `json:"password" binding:"required"`
}

// POST /register
func RegNewCustomerHandler(c *gin.Context) {
	var input RegNewCustomerInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверные данные: " + err.Error()})
		return
	}

	// Проверка на уникальность логина
	var exists models.AuthCredential
	if err := db.DB.Where("login = ?", input.Login).First(&exists).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Логин уже используется"})
		return
	}

	// Хешируем пароль
	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при хешировании пароля"})
		return
	}

	// Создаём user_info
	userInfo := models.UserInfo{
		Phone:    input.Phone,
		FName:    input.FName,
		BirthDay: input.BirthDay,
	}
	if err := db.DB.Create(&userInfo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при создании профиля"})
		return
	}

	// Создаём auth
	auth := models.AuthCredential{
		Login:        input.Login,
		PasswordHash: hash,
	}
	if err := db.DB.Create(&auth).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при создании учетных данных"})
		return
	}

	// Создаём основную запись пользователя
	user := models.User{
		UserType:   models.Customer,
		AuthID:     auth.ID,
		UserInfoID: userInfo.ID,
	}
	if err := db.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при создании пользователя"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Пользователь успешно зарегистрирован", "user_id": user.ID})
}
