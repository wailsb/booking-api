package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"booking-api/internal/domain"

	"github.com/google/uuid"
)

type bookingUseCase struct {
	bookingRepo domain.BookingRepository
	auditRepo   domain.AuditRepository
}

func NewBookingUseCase(bRepo domain.BookingRepository, aRepo domain.AuditRepository) BookingUseCase {
	return &bookingUseCase{
		bookingRepo: bRepo,
		auditRepo:   aRepo,
	}
}

func (uc *bookingUseCase) CreateBooking(ctx context.Context, booking *domain.Booking) (*domain.Booking, error) {
	// 1. Business Validation: Start time must be before End time
	if !booking.StartTime.Before(booking.EndTime) {
		return nil, domain.ErrInvalidTimeRange
	}

	// 2. Business Validation: Cannot book in the past
	if booking.StartTime.Before(time.Now()) {
		return nil, domain.ErrBookingInPast
	}

	// 3. Ensure Bookable resource exists
	_, err := uc.bookingRepo.GetBookableByID(ctx, booking.BookableID)
	if err != nil {
		return nil, fmt.Errorf("bookable check failed: %w", err)
	}

	booking.Status = domain.StatusConfirmed

	// 4. Save to Repository (PostgreSQL handles race conditions via EXCLUDE constraint)
	if err := uc.bookingRepo.CreateBooking(ctx, booking); err != nil {
		return nil, err
	}

	// 5. Audit Event Recording
	auditLog := &domain.AuditLog{
		ID:         uuid.New(),
		Action:     "BOOKING_CREATED",
		EntityType: "BOOKING",
		EntityID:   booking.ID,
		NewState: map[string]interface{}{
			"bookable_id":    booking.BookableID,
			"customer_name":  booking.CustomerName,
			"customer_email": booking.CustomerEmail,
			"start_time":     booking.StartTime,
			"end_time":       booking.EndTime,
			"status":         booking.Status,
		},
		CreatedAt: time.Now(),
	}

	// Record audit log (Actor ID will be attached inside repo from ctx context)
	_ = uc.auditRepo.RecordLog(ctx, auditLog)

	return booking, nil
}

func (uc *bookingUseCase) CancelBooking(ctx context.Context, bookingID uuid.UUID) error {
	// 1. Fetch current booking state
	existing, err := uc.bookingRepo.GetBookingByID(ctx, bookingID)
	if err != nil {
		return err
	}

	if existing.Status == domain.StatusCancelled {
		return errors.New("booking is already cancelled")
	}

	// 2. Extract System User ID performing the action from context
	var actorID uuid.UUID
	if val := ctx.Value("userID"); val != nil {
		if id, ok := val.(uuid.UUID); ok {
			actorID = id
		}
	}

	// 3. Update status in DB
	err = uc.bookingRepo.UpdateBookingStatus(ctx, bookingID, domain.StatusCancelled, actorID)
	if err != nil {
		return err
	}

	// 4. Record Audit Log
	auditLog := &domain.AuditLog{
		ID:         uuid.New(),
		Action:     "BOOKING_CANCELLED",
		EntityType: "BOOKING",
		EntityID:   bookingID,
		OldState:   map[string]interface{}{"status": existing.Status},
		NewState:   map[string]interface{}{"status": domain.StatusCancelled},
		CreatedAt:  time.Now(),
	}

	return uc.auditRepo.RecordLog(ctx, auditLog)
}

func (uc *bookingUseCase) GetBookingByID(ctx context.Context, id uuid.UUID) (*domain.Booking, error) {
	return uc.bookingRepo.GetBookingByID(ctx, id)
}

func (uc *bookingUseCase) ListBookingsByBookable(ctx context.Context, bookableID uuid.UUID, start, end time.Time) ([]domain.Booking, error) {
	return uc.bookingRepo.GetBookingsByBookable(ctx, bookableID, start, end)
}
