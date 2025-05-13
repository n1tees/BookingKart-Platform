package main

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/n1tees/BookingKart-Platform/config"
	"github.com/n1tees/BookingKart-Platform/internal/db"
	"github.com/n1tees/BookingKart-Platform/internal/middleware"
	"github.com/n1tees/BookingKart-Platform/internal/routes"
	"github.com/n1tees/BookingKart-Platform/internal/services"

	_ "github.com/n1tees/BookingKart-Platform/cmd/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	config.LoadEnv()

	time.Sleep(5 * time.Second)

	db.InitDB()
	defer db.CloseDB()

	// запуск фонового обновления статусов
	go func() {
		for {
			services.CheckAndUpdateStatuses()
			time.Sleep(1 * time.Minute)
		}
	}()

	r := gin.Default()
	// Public routes
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	routes.InitAuthRoutes(r)

	// Protected routes
	api := r.Group("/api")
	api.Use(middleware.AuthRequired())

	routes.InitRaceRoutes(api)
	routes.InitUserRoutes(api)
	routes.InitTrackRoutes(api)
	routes.InitKartodromRoutes(api)
	routes.InitPaymentRoutes(api)
	routes.InitBookingRoutes(api)
	routes.InitKartBookingRoutes(api)
	routes.InitKartRoutes(api)

	r.Run() // запускает на :8080
}
