package customer

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/n1tees/BookingKart-Platform/internal/db"
	"github.com/n1tees/BookingKart-Platform/internal/models"
)

type UpdateUserInfoInput struct {
	FName    *string  `json:"fname"`
	SName    *string  `json:"sname"`
	Phone    *string  `json:"phone"`
	Email    *string  `json:"email"`
	Weight   *float64 `json:"weight"`
	Height   *float64 `json:"height"`
	BirthDay *string  `json:"birthday"`
}

func UpdateUserInfoHandler(c *gin.Context) {
	userID := c.GetUint("user_id")

	// ищем пользователя в бд по id
	var user models.User
	if err := db.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Пользователь не найден"})
		return
	}

	var input UpdateUserInfoInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// создаем мапу с обноваленной инфой по полям
	updates := map[string]interface{}{}
	if input.FName != nil {
		updates["f_name"] = *input.FName
	}
	if input.SName != nil {
		updates["s_name"] = *input.SName
	}
	if input.Phone != nil {
		updates["phone"] = *input.Phone
	}
	if input.Email != nil {
		updates["email"] = *input.Email
	}
	if input.Weight != nil {
		updates["weight"] = *input.Weight
	}
	if input.Height != nil {
		updates["height"] = *input.Height
	}
	if input.BirthDay != nil {
		updates["birth_day"] = *input.BirthDay
	}

	// проверяем наличие userinfo по этому id
	if err := db.DB.Model(&models.UserInfo{}).
		Where("id = ?", user.UserInfoID).
		Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при обновлении информации"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Информация обновлена", "updated": updates})
}
