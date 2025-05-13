package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/n1tees/BookingKart-Platform/internal/services"
)

// =================== KART =======================

func GetAvailableKartsHandler(c *gin.Context) {
	kartodromID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	karts, err := services.GetAvailableKarts(uint(kartodromID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, karts)
}

func BookKartHandler(c *gin.Context) {
	kartID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := services.BookKart(uint(kartID)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "карт забронирован"})
}

func FreeKartHandler(c *gin.Context) {
	kartID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := services.FreeKart(uint(kartID)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "карт освобожден"})
}

func SetKartBrokenHandler(c *gin.Context) {
	kartID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := services.SetKartBroken(uint(kartID)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "карт помечен как сломанный"})
}

func RepairKartHandler(c *gin.Context) {
	kartID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := services.RepairKart(uint(kartID)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "карт восстановлен"})
}

func SetKartInStopListHandler(c *gin.Context) {
	kartID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := services.SetKartInStopList(uint(kartID)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "карт добавлен в стоп-лист"})
}
