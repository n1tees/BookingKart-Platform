package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/n1tees/BookingKart-Platform/internal/handlers"
)

func InitAuthRoutes(r *gin.Engine) {
	r.POST("/api/register", handlers.RegisterHandler)
	r.POST("/api/login", handlers.LoginHandler)
}
