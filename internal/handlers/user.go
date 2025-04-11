package handlers

import (
	"net/http"
	"time"

	"github.com/n1tees/BookingKart-Platform/internal/db"
	"github.com/n1tees/BookingKart-Platform/internal/models"

	"github.com/gin-gonic/gin"
)

// Получить всех пользователей (с UserInfo)
func GetAllUsers(c *gin.Context) {
	var users []models.User
	if err := db.DB.Preload("UserInfo").Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении пользователей"})
		return
	}
	c.JSON(http.StatusOK, users)
}

// Получить пользователя по ID
func GetUserByID(c *gin.Context) {
	id := c.Param("id")
	var user models.User
	if err := db.DB.Preload("UserInfo").First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Пользователь не найден"})
		return
	}
	c.JSON(http.StatusOK, user)
}

// Создание пользователя (простой способ — уже созданные AuthID и UserInfoID)
func CreateUser(c *gin.Context) {
	type Input struct {
		FName    string    `json:"fname"`
		Phone    string    `json:"phone"`
		BirthDay time.Time `json:"birthday"`
		Password string    `json:"password"`
	}

	var input Input
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 1. Создаём AuthCredential (логин = phone)
	auth := models.AuthCredential{
		Login:            input.Phone,
		PassWordHashbyte: []byte(input.Password), // на проде — хешировать
	}
	if err := db.DB.Create(&auth).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка создания логина"})
		return
	}

	// 2. Создаём UserInfo
	info := models.UserInfo{
		FName:     input.FName,
		Phone:     input.Phone,
		BirthDay:  input.BirthDay,
		CreatedAt: time.Now(),
	}
	if err := db.DB.Create(&info).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка создания информации"})
		return
	}

	// 3. Создаём пользователя (user_type = customer)
	user := models.User{
		UserType:   "customer",
		AuthID:     auth.ID,
		UserInfoID: info.ID,
	}
	if err := db.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка создания пользователя"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Пользователь создан", "id": user.ID})
}

// Частичное обновление пользователя
func PatchUser(c *gin.Context) {
	id := c.Param("id")
	var user models.User

	if err := db.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Пользователь не найден"})
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.DB.Model(&user).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при обновлении"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Пользователь обновлён", "updated": updates})
}

// Удаление пользователя (+ Auth + Info)
func DeleteUser(c *gin.Context) {
	id := c.Param("id")
	var user models.User

	if err := db.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Пользователь не найден"})
		return
	}

	db.DB.Delete(&models.AuthCredential{}, user.AuthID)
	db.DB.Delete(&models.UserInfo{}, user.UserInfoID)
	db.DB.Delete(&user)

	c.JSON(http.StatusOK, gin.H{"message": "Пользователь удалён"})
}
