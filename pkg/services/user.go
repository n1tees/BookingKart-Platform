package services

import (
	"github.com/n1tees/BookingKart-Platform/pkg/db"
	"github.com/n1tees/BookingKart-Platform/pkg/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"errors"
	"fmt"
	"strings"
	"sync"
)

// Стуктура ввода для смены пароля
type ChangePasswordInput struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
	RepeatNew   string `json:"repeat_new_password" binding:"required"`
}

// Частичное обновление пользователя
func UpdateProfile(userID uint, updates map[string]interface{}) error {

	user, err := searchUserByID(userID)
	if err != nil {
		return err
	}

	filtered := make(map[string]interface{})
	for key, value := range updates {
		if value == nil {
			continue
		}
		switch v := value.(type) {
		case string:
			if strings.TrimSpace(v) == "" {
				continue
			}
			filtered[key] = v
		case float64:
			if v == 0 {
				continue
			}
			filtered[key] = v
		case int, int32, int64:
			if fmt.Sprintf("%v", v) == "0" {
				continue
			}
			filtered[key] = v
		case bool:
			if !v {
				continue
			}
			filtered[key] = v
		default:
			filtered[key] = v
		}
	}

	if len(filtered) == 0 {
		return errors.New("пустой запрос")
	}

	if err := db.DB.Model(&models.Profile{}).Where("id = ?", user.ProfileID).Updates(filtered).Error; err != nil {
		return errors.New("ошибка при обновлении профиля")
	}

	return nil
}

// Зайти в свой профиль
func GetUserInfo(userID uint) (*models.Profile, error) {

	user, err := searchUserByID(userID)
	if err != nil {
		return nil, err
	}

	var profile models.Profile
	if err := db.DB.First(&profile, user.ProfileID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("пользователь не найден")

		} else {
			return nil, errors.New("ошибка при поиске пользователя")
		}
	}

	return &profile, nil
}

// Изменить пароль
func ChangePassword(userID uint, input ChangePasswordInput) error {

	var mu sync.Mutex
	mu.Lock()
	defer mu.Unlock()

	// Находим пользователя
	user, err := searchUserByID(userID)
	if err != nil {
		return err
	}

	// Вытаскиваем auth
	var auth models.AuthCredential
	if err := db.DB.First(&auth, user.AuthID).Error; err != nil {
		return errors.New("учётные данные не найдены")
	}

	// Сравниеваем новые пароли
	if input.NewPassword != input.RepeatNew {
		return errors.New("новые пароли не совпадают")
	}

	// Сравниеваем со старым паролем
	if err := bcrypt.CompareHashAndPassword(auth.PasswordHash, []byte(input.OldPassword)); err != nil {
		return errors.New("старый пароль неверен")
	}

	// Генерим хещ
	newHash, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("ошибка при хешировании пароля")
	}

	// Обоновляем пароль
	if err := db.DB.Model(&auth).Update("password_hash", newHash).Error; err != nil {
		return errors.New("ошибка при обновлении пароля")
	}

	return nil
}

// поиск пользователя с заданными auth_credential
func searchUserByAuth(auth *models.AuthCredential) (*models.User, error) {

	var user models.User
	if err := db.DB.Where("auth_id = ?", auth.ID).First(&user).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("пользователь не найден")

		} else {
			return nil, errors.New("ошибка при поиске пользователя")
		}
	}

	return &user, nil
}

// поиск по id
func searchUserByID(id uint) (*models.User, error) {

	var user models.User
	if err := db.DB.First(&user, id).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("профиль не найден")

		} else {
			return nil, errors.New("ошибка при поиске профиля")
		}
	}
	return &user, nil
}

// поиск по логину
func searchAuthByLogin(login string) (*models.AuthCredential, error) {

	var auth models.AuthCredential
	err := db.DB.Where("login = ?", login).First(&auth).Error
	if err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("учётные данные не найдены")

		} else {
			return nil, errors.New("ошибка при поиске учетных данных")
		}
	}

	return &auth, nil
}

// поиск по телефону
func searchProfileByPhone(phone string) (*models.Profile, error) {

	var profile models.Profile
	if err := db.DB.Where("phone = ?", phone).First(&profile).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("профиль не найден")

		} else {
			return nil, errors.New("ошибка при поиске профиля")
		}
	}
	return &profile, nil
}
