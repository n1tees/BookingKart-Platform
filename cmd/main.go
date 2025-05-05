package main

import (
	"github.com/n1tees/BookingKart-Platform/config"
	_ "github.com/n1tees/BookingKart-Platform/internal/db"
)

func main() {
	config.LoadEnv()
	//db.InitDB()

	// Здесь позже запустим сервер
}
