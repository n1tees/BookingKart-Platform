package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/n1tees/BookingKart-Platform/pkg/services"
)

// GetUserHandler godoc
// @Summary Получить профиль пользователя
// @Tags user
// @Param id path int true "ID пользователя"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /user/{id} [get]
func GetUserHandler(c *gin.Context) {
	userID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	profile, err := services.GetUserInfo(uint(userID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, profile)
}

type UpdateUserRequest map[string]interface{}

// UpdateUserHandler godoc
// @Summary Обновить профиль пользователя
// @Tags user
// @Param id path int true "ID пользователя"
// @Param updates body handlers.UpdateUserRequest true "Обновляемые данные"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /user/{id} [patch]
func UpdateUserHandler(c *gin.Context) {
	userID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var updates UpdateUserRequest
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный формат запроса"})
		return
	}

	if err := services.UpdateProfile(uint(userID), updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "профиль обновлен"})
}

// ChangePasswordHandler godoc
// @Summary Сменить пароль пользователя
// @Tags user
// @Param id path int true "ID пользователя"
// @Param input body services.ChangePasswordInput true "Данные для смены пароля"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /user/{id}/change-password [post]
func ChangePasswordHandler(c *gin.Context) {
	userID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var input services.ChangePasswordInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный формат запроса"})
		return
	}

	if err := services.ChangePassword(uint(userID), input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "пароль успешно изменен"})
}
