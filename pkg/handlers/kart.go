package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/n1tees/BookingKart-Platform/pkg/services"
)

// GetAvailableKartsHandler godoc
// @Summary Получить доступные карты в картодроме
// @Tags kart
// @Param id path int true "ID картодрома"
// @Success 200 {array} models.Kart
// @Failure 500 {object} map[string]string
// @Router /kartodrom/{id}/free-karts [get]
func GetAvailableKartsHandler(c *gin.Context) {
	kartodromID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	karts, err := services.GetAvailableKarts(uint(kartodromID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, karts)
}

// BookKartHandler godoc
// @Summary Забронировать карт
// @Tags kart
// @Param id path int true "ID карта"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /karts/{id}/book [post]
func BookKartHandler(c *gin.Context) {
	kartID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := services.BookKart(uint(kartID)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "карт забронирован"})
}

// FreeKartHandler godoc
// @Summary Освободить карт
// @Tags kart
// @Param id path int true "ID карта"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /karts/{id}/free [post]
func FreeKartHandler(c *gin.Context) {
	kartID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := services.FreeKart(uint(kartID)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "карт освобожден"})
}

// SetKartBrokenHandler godoc
// @Summary Пометить карт как сломанный
// @Tags kart
// @Param id path int true "ID карта"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /karts/{id}/broken [post]
func SetKartBrokenHandler(c *gin.Context) {
	kartID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := services.SetKartBroken(uint(kartID)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "карт помечен как сломанный"})
}

// RepairKartHandler godoc
// @Summary Восстановить карт
// @Tags kart
// @Param id path int true "ID карта"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /karts/{id}/repair [post]
func RepairKartHandler(c *gin.Context) {
	kartID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := services.RepairKart(uint(kartID)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "карт восстановлен"})
}

// SetKartInStopListHandler godoc
// @Summary Добавить карт в стоп-лист
// @Tags kart
// @Param id path int true "ID карта"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /karts/{id}/stoplist [post]
func SetKartInStopListHandler(c *gin.Context) {
	kartID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := services.SetKartInStopList(uint(kartID)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "карт добавлен в стоп-лист"})
}
