package main

import (
	"time"

	_ "github.com/gin-gonic/gin"
	"github.com/n1tees/BookingKart-Platform/config"
	"github.com/n1tees/BookingKart-Platform/internal/db"
	_ "github.com/n1tees/BookingKart-Platform/internal/routes"
	"github.com/n1tees/BookingKart-Platform/internal/services"
)

func main() {
	config.LoadEnv()

	time.Sleep(5 * time.Second)

	db.InitDB()
	defer db.CloseDB()

	//var res *[]models.Payment
	//var res float64
	var err error
	err = services.RefillMyBalance(1, 1200.0)
	if err != nil {
		print(err.Error())
		// } else {
		// 	fmt.Printf("Баланс пользователя: %.2f\n", res)
		// r := gin.Default()
		// routes.InitAuthRoutes(r)
		// r.Run() // запускает на :8080

	}
}
