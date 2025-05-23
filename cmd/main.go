package main

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/n1tees/BookingKart-Platform/config"
	"github.com/n1tees/BookingKart-Platform/pkg/db"
	"github.com/n1tees/BookingKart-Platform/pkg/middleware"
	"github.com/n1tees/BookingKart-Platform/pkg/routes"
	"github.com/n1tees/BookingKart-Platform/pkg/services"

	_ "github.com/n1tees/BookingKart-Platform/docs"
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
	// r.Use(cors.New(cors.Config{
	// 	AllowOrigins:     []string{"http://localhost:5173"},
	// 	AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
	// 	AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
	// 	ExposeHeaders:    []string{"Content-Length"},
	// 	AllowCredentials: true,
	// 	MaxAge:           12 * time.Hour,
	// }))

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
