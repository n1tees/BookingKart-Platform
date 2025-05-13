package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/n1tees/BookingKart-Platform/internal/models"
	"github.com/n1tees/BookingKart-Platform/internal/services"
)

type CreateRaceRequest struct {
	TrackID   uint            `json:"track_id"`
	Date      string          `json:"date"`
	TimeStart string          `json:"time_start"`
	RaceType  models.RaceType `json:"race_type"`
	Laps      uint            `json:"laps"`
	Duration  uint            `json:"duration"`
}

// CreateRaceHandler godoc
// @Summary Создать гонку
// @Tags race
// @Param input body handlers.CreateRaceRequest true "Данные гонки"
// @Success 201 {object} map[string]uint
// @Failure 400 {object} map[string]string
// @Router /races [post]
func CreateRaceHandler(c *gin.Context) {
	var req CreateRaceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный формат запроса"})
		return
	}

	parsedDate, err := MakeDateByString(req.Date)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный формат даты (YYYY-MM-DD)"})
		return
	}
	parsedTime, err := MakeTimeByString(req.TimeStart)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный формат времени (HH:MM)"})
		return
	}

	raceInput := services.RaceInput{
		TrackID:         req.TrackID,
		Date:            req.Date,
		TimeStart:       req.TimeStart,
		RaceType:        req.RaceType,
		Laps:            req.Laps,
		Duration:        req.Duration,
		ParsedDate:      parsedDate,
		ParsedTimeStart: parsedTime,
	}

	raceID, err := services.CreateRace(raceInput)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"race_id": raceID})
}

// StartRaceHandler godoc
// @Summary Стартовать гонку
// @Tags race
// @Param id path int true "ID гонки"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /races/{id}/start [post]
func StartRaceHandler(c *gin.Context) {
	raceID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := services.StartRace(uint(raceID)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "гонка стартовала"})
}

// FinishRaceHandler godoc
// @Summary Завершить гонку
// @Tags race
// @Param id path int true "ID гонки"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /races/{id}/finish [post]
func FinishRaceHandler(c *gin.Context) {
	raceID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := services.FinishRace(uint(raceID)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "гонка завершена"})
}

// CancelRaceHandler godoc
// @Summary Отменить гонку
// @Tags race
// @Param id path int true "ID гонки"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /races/{id}/cancel [post]
func CancelRaceHandler(c *gin.Context) {
	raceID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := services.CancelRace(uint(raceID)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "гонка отменена"})
}

// =================== RACERIDER =======================

type RegisterRiderRequest struct {
	RiderID      uint `json:"rider_id"`
	ResultTypeID uint `json:"result_type_id"`
}

// RegisterRiderHandler godoc
// @Summary Зарегистрировать участника в гонке
// @Tags race
// @Param id path int true "ID гонки"
// @Param input body handlers.RegisterRiderRequest true "Данные участника"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /races/{id}/riders [post]
func RegisterRiderHandler(c *gin.Context) {
	raceID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var req RegisterRiderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный формат запроса"})
		return
	}

	if err := services.RegisterRider(uint(raceID), req.RiderID, req.ResultTypeID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "участник добавлен в гонку"})
}

// RemoveRiderHandler godoc
// @Summary Удалить участника из гонки
// @Tags race
// @Param id path int true "ID гонки"
// @Param riderId path int true "ID участника"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /races/{id}/riders/{riderId} [delete]
func RemoveRiderHandler(c *gin.Context) {
	raceID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	riderID, _ := strconv.ParseUint(c.Param("riderId"), 10, 64)
	if err := services.RemoveRider(uint(raceID), uint(riderID)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "участник удален из гонки"})
}

type AddResultRequest struct {
	RiderID        uint `json:"rider_id"`
	ResultTypeID   uint `json:"result_type_id"`
	PersonalResult uint `json:"personal_result"`
}

// AddRaceResultHandler godoc
// @Summary Добавить результат участника
// @Tags race
// @Param id path int true "ID гонки"
// @Param input body handlers.AddResultRequest true "Данные результата"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /races/{id}/results [post]
func AddRaceResultHandler(c *gin.Context) {
	raceID, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var req AddResultRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный формат запроса"})
		return
	}

	if err := services.AddRaceResult(uint(raceID), req.RiderID, req.ResultTypeID, req.PersonalResult); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "результат участника добавлен"})
}
