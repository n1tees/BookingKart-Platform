package main

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/n1tees/BookingKart-Platform/config"
	"github.com/n1tees/BookingKart-Platform/internal/db"
	"github.com/n1tees/BookingKart-Platform/internal/routes"
)

func main() {
	config.LoadEnv()

	time.Sleep(5 * time.Second)

	db.InitDB()
	defer db.CloseDB()

	r := gin.Default()
	routes.InitAuthRoutes(r)
	r.Run() // запускает на :8080
}
