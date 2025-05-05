package main

import (
	"time"

	"github.com/n1tees/BookingKart-Platform/config"
	"github.com/n1tees/BookingKart-Platform/internal/db"
)

func main() {
	config.LoadEnv()

	time.Sleep(2 * time.Second)

	db.InitDB()

	// Здесь позже запустим сервер
}
