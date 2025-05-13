package services

import (
	"fmt"
	"time"

	"github.com/n1tees/BookingKart-Platform/pkg/db"
	"github.com/n1tees/BookingKart-Platform/pkg/models"
)

func CheckAndUpdateStatuses() {
	fmt.Println("[Scheduler] Запуск проверки и обновления статусов...")

	updateRaces()
	updateBookings()
	updateKartBookings()
}

func updateRaces() {
	now := time.Now()

	var races []models.Race
	db.DB.Find(&races)

	for _, race := range races {
		startTime := time.Date(race.Date.Year(), race.Date.Month(), race.Date.Day(),
			race.TimeStart.Hour(), race.TimeStart.Minute(), race.TimeStart.Second(), 0, race.TimeStart.Location())

		if race.Status == models.RaceCreate && now.After(startTime) {
			db.DB.Model(&race).Update("status", models.RaceStart)
		}

		endTime := startTime.Add(time.Duration(race.Duration) * time.Minute)
		if race.Status == models.RaceStart && now.After(endTime) {
			db.DB.Model(&race).Update("status", models.RaceFinish)
		}
	}
}

func updateBookings() {
	now := time.Now()

	var bookings []models.Booking
	db.DB.Find(&bookings)

	for _, booking := range bookings {
		start := booking.StartTime.Time
		end := booking.EndTime.Time

		if booking.Status == models.BookingReserve && now.After(start) {
			db.DB.Model(&booking).Update("status", models.BookingActive)
		}
		if booking.Status == models.BookingActive && now.After(end) {
			db.DB.Model(&booking).Update("status", models.BookingClose)
		}
	}
}

func updateKartBookings() {
	var kartBookings []models.KartBooking
	db.DB.Preload("Booking").Preload("Kart").Find(&kartBookings)

	for _, kb := range kartBookings {
		if kb.Booking.Status == models.BookingClose || kb.Booking.Status == models.BookingCancel {
			if kb.Kart.Status != models.Available {
				db.DB.Model(&kb.Kart).Update("status", models.Available)
			}
		}
	}
}
