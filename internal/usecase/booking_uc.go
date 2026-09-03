package usecase

import (
	"context"
	"time"

	"booking-api/internal/domain"

	"github.com/google/uuid"
)

type BookingUseCase interface {
	CreateBooking(ctx context.Context, booking *domain.Booking) (*domain.Booking, error)
	CancelBooking(ctx context.Context, bookingID uuid.UUID) error
	GetBookingByID(ctx context.Context, id uuid.UUID) (*domain.Booking, error)
	ListBookingsByBookable(ctx context.Context, bookableID uuid.UUID, start, end time.Time) ([]domain.Booking, error)
}
