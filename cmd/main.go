package main

import (
	"log"
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

	// var input = map[string]interface{}{
	// 	"FName": "Nikitius",
	// 	"SName": " ",
	// 	"Email": "sadfsdfsdfasdf@gmail.com",
	// }

	//var track *[]models.Kart
	var err error

	err = services.SetKartInStopList(1)
	if err != nil {
		log.Printf(err.Error())
	}
	//log.Printf("%+v", track)

	// r := gin.Default()
	// routes.InitAuthRoutes(r)
	// r.Run() // запускает на :8080
}
