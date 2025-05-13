package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/n1tees/BookingKart-Platform/internal/services"
)

func GetAvailableTracksHandler(c *gin.Context) {
	kartodromID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	tracks, err := services.GetAvailableTracks(uint(kartodromID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tracks)
}

func GetTrackByIDHandler(c *gin.Context) {
	trackID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	track, err := services.GetTrackByID(uint(trackID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, track)
}
