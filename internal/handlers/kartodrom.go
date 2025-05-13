package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/n1tees/BookingKart-Platform/internal/services"
)

func GetKartodromsHandler(c *gin.Context) {
	city := c.Query("city")
	var cityPtr *string
	if city != "" {
		cityPtr = &city
	}

	kartodroms, err := services.GetKartodromsList(cityPtr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, kartodroms)
}
