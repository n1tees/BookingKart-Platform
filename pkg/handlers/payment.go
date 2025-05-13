package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/n1tees/BookingKart-Platform/pkg/services"
)

// GetPaymentsHandler godoc
// @Summary Получить историю платежей пользователя
// @Tags payment
// @Param id path int true "ID пользователя"
// @Success 200 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /user/{id}/payments [get]
func GetPaymentsHandler(c *gin.Context) {
	userID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	payments, err := services.GetMyPayments(uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, payments)
}

// GetBalanceHandler godoc
// @Summary Получить текущий баланс пользователя
// @Tags payment
// @Param id path int true "ID пользователя"
// @Success 200 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /user/{id}/balance [get]
func GetBalanceHandler(c *gin.Context) {
	userID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	balance, err := services.GetBalance(uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"balance": balance})
}

type AmountRequest struct {
	Amount float64 `json:"amount"`
}

// RefillBalanceHandler godoc
// @Summary Пополнить баланс пользователя
// @Tags payment
// @Param id path int true "ID пользователя"
// @Param amount body handlers.AmountRequest true "Сумма пополнения"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /user/{id}/refill [post]
func RefillBalanceHandler(c *gin.Context) {
	userID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var req AmountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный формат запроса"})
		return
	}

	if err := services.RefillMyBalance(uint(userID), req.Amount); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "баланс успешно пополнен"})
}

// RefundBalanceHandler godoc
// @Summary Сделать возврат средств пользователю
// @Tags payment
// @Param id path int true "ID пользователя"
// @Param amount body handlers.AmountRequest true "Сумма возврата"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /user/{id}/refund [post]
func RefundBalanceHandler(c *gin.Context) {
	userID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var req AmountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный формат запроса"})
		return
	}

	if err := services.RefundToBalance(uint(userID), req.Amount); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "возврат выполнен успешно"})
}
