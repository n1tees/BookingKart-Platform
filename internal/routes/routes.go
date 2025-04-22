package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/n1tees/BookingKart-Platform/internal/handlers/auth"
	handlers "github.com/n1tees/BookingKart-Platform/internal/handlers/temp"
)

func InitRoutes(router *gin.Engine) {
	// Пример
	router.POST("/auth/login", auth.LoginHandler)

	karts := router.Group("/karts")
	{
		karts.GET("", handlers.GetAllKarts)
		karts.POST("", handlers.CreateKart)
	}
}
