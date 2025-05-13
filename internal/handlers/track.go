package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/n1tees/BookingKart-Platform/internal/services"
)

// GetAvailableTracksHandler godoc
// @Summary Получить доступные треки картодрома
// @Tags track
// @Param id path int true "ID картодрома"
// @Success 200 {array} models.Track
// @Failure 404 {object} map[string]string
// @Router /kartodrom/{id}/tracks [get]
func GetAvailableTracksHandler(c *gin.Context) {
	kartodromID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	tracks, err := services.GetAvailableTracks(uint(kartodromID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tracks)
}

// GetTrackByIDHandler godoc
// @Summary Получить трек по ID
// @Tags track
// @Param id path int true "ID трека"
// @Success 200 {object} models.Track
// @Failure 404 {object} map[string]string
// @Router /track/{id} [get]
func GetTrackByIDHandler(c *gin.Context) {
	trackID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	track, err := services.GetTrackByID(uint(trackID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, track)
}
