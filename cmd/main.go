package main

import (
	"errors"
	"fmt"
	"time"

	_ "github.com/gin-gonic/gin"
	"github.com/n1tees/BookingKart-Platform/config"
	"github.com/n1tees/BookingKart-Platform/internal/db"
	"github.com/n1tees/BookingKart-Platform/internal/models"
	_ "github.com/n1tees/BookingKart-Platform/internal/routes"
	"gorm.io/gorm"
)

func main() {
	config.LoadEnv()

	time.Sleep(5 * time.Second)

	db.InitDB()
	defer db.CloseDB()

	//var result uint
	//var err error

	// err = services.FinishRace(2)
	// if err != nil {
	// 	log.Printf(err.Error())
	// }

	var race models.Race
	if err := db.DB.First(&race, 1).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			fmt.Println("заезд не найден")

		} else {
			fmt.Println(err.Error())
			return
		}
	}

	//log.Printf("%+v", result)

	// r := gin.Default()
	// routes.InitAuthRoutes(r)
	// r.Run() // запускает на :8080
}
