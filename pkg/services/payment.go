package services

import (
	"errors"
	"sync"
	"time"

	"github.com/n1tees/BookingKart-Platform/pkg/db"
	"github.com/n1tees/BookingKart-Platform/pkg/models"
	"gorm.io/gorm"
)

// Получить история пополнения баланса
func GetMyPayments(userID uint) (*[]models.Payment, error) {
	var payments []models.Payment

	if err := db.DB.Where("user_id = ?", userID).Find(&payments).Error; err != nil {
		return nil, errors.New("ошибка при получении списка платежей")
	}

	return &payments, nil
}

// Пополнить баланс
func RefillMyBalance(userID uint, amount float64) error {
	if amount <= 0 {
		return errors.New("сумма пополнения должна быть больше нуля")
	}

	// Блокировка
	var mu sync.Mutex
	mu.Lock()
	defer mu.Unlock()

	// Транзакция
	err := db.DB.Transaction(func(tx *gorm.DB) error {
		// Находим пользователя
		var user models.User
		if err := tx.First(&user, userID).Error; err != nil {
			return errors.New("пользователь не найден")
		}

		// Находим профиль
		var profile models.Profile
		if err := tx.First(&profile, user.ProfileID).Error; err != nil {
			return errors.New("профиль не найден")
		}

		// Симулируем оплату через банк
		if !simulateBankPayment(userID, amount) {
			return errors.New("платёж не подтвержден банком")
		}

		// Увеличиваем баланс
		profile.Balance += amount
		if err := tx.Save(&profile).Error; err != nil {
			return errors.New("ошибка при обновлении баланса")
		}

		// Записываем пополнение в payments
		payment := models.Payment{
			UserID: userID,
			Amount: amount,
			Date:   time.Now(),
		}
		if err := tx.Create(&payment).Error; err != nil {
			return errors.New("ошибка при записи платежа")
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

// Возврат средств
func RefundToBalance(userID uint, amount float64) error {
	if amount <= 0 {
		return errors.New("некорректная сумма возврата")
	}

	return db.DB.Transaction(func(tx *gorm.DB) error {
		var user models.User
		if err := tx.First(&user, userID).Error; err != nil {
			return errors.New("пользователь не найден")
		}

		var profile models.Profile
		if err := tx.First(&profile, user.ProfileID).Error; err != nil {
			return errors.New("профиль пользователя не найден")
		}

		profile.Balance += amount
		if err := tx.Save(&profile).Error; err != nil {
			return errors.New("ошибка при возврате денег на баланс")
		}

		// Записываем пополнение в payments
		payment := models.Payment{
			UserID: userID,
			Amount: -amount,
			Date:   time.Now(),
		}
		if err := tx.Create(&payment).Error; err != nil {
			return errors.New("ошибка при записи платежа")
		}

		return nil
	})
}

// Получить текущий баланс
func GetBalance(userID uint) (float64, error) {
	var user models.User
	if err := db.DB.First(&user, userID).Error; err != nil {
		return 0, errors.New("пользователь не найден")
	}

	var profile models.Profile
	if err := db.DB.First(&profile, user.ProfileID).Error; err != nil {
		return 0, errors.New("профиль пользователя не найден")
	}

	return profile.Balance, nil
}

// Списать деньги с баланса
func ChargeFromBalance(userID uint, amount float64) error {
	if amount <= 0 {
		return errors.New("некорректная сумма списания")
	}

	return db.DB.Transaction(func(tx *gorm.DB) error {
		var user models.User
		if err := tx.First(&user, userID).Error; err != nil {
			return errors.New("пользователь не найден")
		}

		var profile models.Profile
		if err := tx.First(&profile, user.ProfileID).Error; err != nil {
			return errors.New("профиль пользователя не найден")
		}

		if profile.Balance < amount {
			return errors.New("недостаточно средств на балансе")
		}

		profile.Balance -= amount
		if err := tx.Save(&profile).Error; err != nil {
			return errors.New("ошибка при списании средств с баланса")
		}

		return nil
	})
}

// симуляция проверки банковского перевода
func simulateBankPayment(userID uint, value float64) bool {

	userID += 1
	value += 1.0
	return true
}
