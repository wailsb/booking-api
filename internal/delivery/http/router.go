package http

import (
	"net/http"

	"booking-api/internal/delivery/http/middleware"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
)

func NewRouter(bookingHandler *BookingHandler, jwtSecret []byte) http.Handler {
	r := chi.NewRouter()

	// Global Middlewares
	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.Recoverer)

	// API Routes Group
	r.Route("/api/v1", func(r chi.Router) {

		// Public Routes (e.g., viewing available slots or public endpoints)
		r.Group(func(r chi.Router) {
			r.Get("/bookables/{bookableID}/bookings", bookingHandler.ListBookingsByBookable)
		})

		// Protected Routes (Requires authenticated JWT token)
		r.Group(func(r chi.Router) {
			r.Use(middleware.AuthMiddleware(jwtSecret)) // Require valid token for any role

			r.Post("/bookings", bookingHandler.CreateBooking)
			r.Get("/bookings/{id}", bookingHandler.GetBookingByID)
			r.Post("/bookings/{id}/cancel", bookingHandler.CancelBooking)
		})

		// Admin-Only Protected Routes Example
		r.Group(func(r chi.Router) {
			r.Use(middleware.AuthMiddleware(jwtSecret, "ADMIN")) // Require "ADMIN" role specifically
			// Place admin-only routes here...
		})
	})

	return r
}
