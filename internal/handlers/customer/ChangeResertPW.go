package customer

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"github.com/n1tees/BookingKart-Platform/internal/db"
	"github.com/n1tees/BookingKart-Platform/internal/models"
)

// стуктура ввода
type ChangePasswordInput struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
	RepeatNew   string `json:"repeat_new_password" binding:"required"`
}

func ChangePasswordHandler(c *gin.Context) {
	userID := c.GetUint("user_id")

	// проверяем наличие такого пользователя
	var user models.User
	if err := db.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Пользователь не найден"})
		return
	}

	// вытаскиваем auth данные
	var auth models.AuthCredential
	if err := db.DB.First(&auth, user.AuthID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Учетные данные не найдены"})
		return
	}

	// вытаскиваем данные из запроса
	var input ChangePasswordInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// проверяем сходство введенного пароля с текущим
	if err := bcrypt.CompareHashAndPassword(auth.PasswordHash, []byte(input.OldPassword)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Старый пароль неверен"})
		return
	}

	// проверяем соответсвтие новых паролей
	if input.NewPassword != input.RepeatNew {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Новые пароли не совпадают"})
		return
	}

	if err := GenerateAndSetHashPW(&auth, input.NewPassword); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при обновлении пароля"})

	}

	// передаем инфу об успешной смене
	c.JSON(http.StatusOK, gin.H{"message": "Пароль успешно обновлён"})
}

// Запрос на сброс пароля по Email
func RequestResetPassword(c *gin.Context) {
	type Request struct {
		Email string `json:"email" binding:"required,email"`
	}

	var input Request
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный формат email"})
		return
	}

	var info models.UserInfo
	if err := db.DB.Where("email = ?", input.Email).First(&info).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Пользователь с таким email не найден"})
		return
	}

	// В демонстрационной версии просто вернём его
	c.JSON(http.StatusOK, gin.H{"message": "Код отправлен на email"})
}

// по поводу сбросу, тут есть так сказать дыра в плане логики
// так как я нигде не храню отправленные коды, и на самом деле их никуда не отправляю
// то в полноый степени механизм восстановления пароля не работает, в будущем если будет время и желание попробую это реалдизовать

// Сброс пароля по email (без ввода старого пароля)
func ResetPasswordByEmail(c *gin.Context) {
	type ResetInput struct {
		Email          string `json:"email" binding:"required,email"`
		NewPassword    string `json:"new_password" binding:"required"`
		RepeatPassword string `json:"repeat_password" binding:"required"`
		Code           int    `json:"code" binding:"required"`
	}

	// вытаскиваем полученные данные
	var input ResetInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// проверка кода
	if input.Code != 123456 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный код подтверждения"})
		return
	}

	// проверка паролей на сходство
	if input.NewPassword != input.RepeatPassword {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Пароли не совпадают"})
		return
	}

	// поиск записи по мылу
	var info models.UserInfo
	if err := db.DB.Where("email = ?", input.Email).First(&info).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Пользователь не найден по email"})
		return
	}

	// поиск пользователя по профилю
	var user models.User
	if err := db.DB.Where("user_info_id = ?", info.ID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Пользователь не найден"})
		return
	}

	// поиск auth по пользователю
	var auth models.AuthCredential
	if err := db.DB.First(&auth, user.AuthID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Учетные данные не найдены"})
		return
	}

	if err := GenerateAndSetHashPW(&auth, input.NewPassword); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при обновлении пароля"})

	}

	c.JSON(http.StatusOK, gin.H{"message": "Пароль успешно сброшен"})
}

func GenerateAndSetHashPW(auth *models.AuthCredential, newPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	return db.DB.Model(auth).Update("password_hash", hash).Error
}
