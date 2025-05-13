package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/n1tees/BookingKart-Platform/pkg/handlers"
)

func InitAuthRoutes(r *gin.Engine) {
	r.POST("/api/register", handlers.RegisterHandler)
	r.POST("/api/login", handlers.LoginHandler)
}

func InitRaceRoutes(r *gin.RouterGroup) {
	r.POST("/api/races", handlers.CreateRaceHandler)
	r.POST("/api/races/:id/start", handlers.StartRaceHandler)
	r.POST("/api/races/:id/finish", handlers.FinishRaceHandler)
	r.POST("/api/races/:id/cancel", handlers.CancelRaceHandler)

	r.POST("/api/races/:id/riders", handlers.RegisterRiderHandler)
	r.DELETE("/api/races/:id/riders/:riderId", handlers.RemoveRiderHandler)

	r.POST("/api/races/:id/results", handlers.AddRaceResultHandler)
}

func InitUserRoutes(r *gin.RouterGroup) {
	r.GET("/api/user/:id", handlers.GetUserHandler)
	r.PATCH("/api/user/:id", handlers.UpdateUserHandler)
	r.POST("/api/user/:id/change-password", handlers.ChangePasswordHandler)
}

func InitTrackRoutes(r *gin.RouterGroup) {
	r.GET("/api/kartodrom/:id/tracks", handlers.GetAvailableTracksHandler)
	r.GET("/api/track/:id", handlers.GetTrackByIDHandler)
}

func InitKartodromRoutes(r *gin.RouterGroup) {
	r.GET("/api/kartodroms", handlers.GetKartodromsHandler)
}

func InitPaymentRoutes(r *gin.RouterGroup) {
	r.GET("/api/user/:id/payments", handlers.GetPaymentsHandler)
	r.GET("/api/user/:id/balance", handlers.GetBalanceHandler)
	r.POST("/api/user/:id/refill", handlers.RefillBalanceHandler)
	r.POST("/api/user/:id/refund", handlers.RefundBalanceHandler)
}

func InitBookingRoutes(r *gin.RouterGroup) {
	r.POST("/api/bookings", handlers.CreateBookingHandler)
	r.POST("/api/bookings/:id/activate", handlers.ActivateBookingHandler)
	r.POST("/api/bookings/:id/close", handlers.CloseBookingHandler)
	r.POST("/api/bookings/:id/cancel", handlers.CancelBookingHandler)
	r.GET("/api/kartodrom/:id/bookings", handlers.GetBookingsByDateHandler)
}

func InitKartBookingRoutes(r *gin.RouterGroup) {
	r.GET("/api/kartodrom/:id/karts", handlers.GetAvailableKartsForBookingHandler)
	r.POST("/api/kartbookings", handlers.ReserveKartHandler)
	r.POST("/api/kartbookings/:bookingId/:kartId/activate", handlers.ActivateKartBookingHandler)
	r.POST("/api/kartbookings/:bookingId/:kartId/finish", handlers.FinishKartBookingHandler)
	r.POST("/api/kartbookings/:bookingId/:kartId/cancel", handlers.CancelKartBookingHandler)
}

func InitKartRoutes(r *gin.RouterGroup) {
	r.GET("/api/kartodrom/:id/free-karts", handlers.GetAvailableKartsHandler)
	r.POST("/api/karts/:id/book", handlers.BookKartHandler)
	r.POST("/api/karts/:id/free", handlers.FreeKartHandler)
	r.POST("/api/karts/:id/broken", handlers.SetKartBrokenHandler)
	r.POST("/api/karts/:id/repair", handlers.RepairKartHandler)
	r.POST("/api/karts/:id/stoplist", handlers.SetKartInStopListHandler)
}
