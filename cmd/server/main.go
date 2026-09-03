package main

import (
	"context"
	"log"
	"net/http"
	"os"

	deliveryHTTP "booking-api/internal/delivery/http"
	"booking-api/internal/repository/postgres"
	"booking-api/internal/usecase"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	jwtSecret := []byte(os.Getenv("JWT_SECRET"))
	dbURL := os.Getenv("DATABASE_URL")

	dbPool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer dbPool.Close()

	// 1. Repositories (Data Access Layer)
	bookingRepo := postgres.NewBookingRepository(dbPool)
	auditRepo := postgres.NewAuditRepository(dbPool)

	// 2. Use Case / Service Layer (Business Logic)
	bookingUC := usecase.NewBookingUseCase(bookingRepo, auditRepo)

	// 3. Handler (HTTP Adapter Layer)
	bookingHandler := deliveryHTTP.NewBookingHandler(bookingUC)

	// 4. Router (Transport & Middleware Wiring)
	router := deliveryHTTP.NewRouter(bookingHandler, jwtSecret)

	log.Println("Server running on port :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
}
