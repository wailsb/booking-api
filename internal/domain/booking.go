package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type BookableType string

const (
	BookableTypeService  BookableType = "SERVICE"
	BookableTypePhysical BookableType = "PHYSICAL"
)

type BookingStatus string

const (
	StatusPending   BookingStatus = "PENDING"
	StatusConfirmed BookingStatus = "CONFIRMED"
	StatusCancelled BookingStatus = "CANCELLED"
	StatusCompleted BookingStatus = "COMPLETED"
)

type Bookable struct {
	ID          uuid.UUID              `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Type        BookableType           `json:"type"`
	Metadata    map[string]interface{} `json:"metadata"`
	CreatedBy   *uuid.UUID             `json:"created_by,omitempty"`
	UpdatedBy   *uuid.UUID             `json:"updated_by,omitempty"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

type Booking struct {
	ID            uuid.UUID     `json:"id"`
	BookableID    uuid.UUID     `json:"bookable_id"`
	CustomerName  string        `json:"customer_name"`
	CustomerEmail string        `json:"customer_email"`
	CustomerPhone string        `json:"customer_phone"`
	StartTime     time.Time     `json:"start_time"`
	EndTime       time.Time     `json:"end_time"`
	Status        BookingStatus `json:"status"`
	CreatedBy     *uuid.UUID    `json:"created_by,omitempty"`
	UpdatedBy     *uuid.UUID    `json:"updated_by,omitempty"`
	CreatedAt     time.Time     `json:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at"`
}

// Repositories
type BookingRepository interface {
	CreateBookable(ctx context.Context, bookable *Bookable) error
	GetBookableByID(ctx context.Context, id uuid.UUID) (*Bookable, error)
	GetBookingByID(ctx context.Context, id uuid.UUID) (*Booking, error)
	ListBookingsByBookable(ctx context.Context, bookableID uuid.UUID, start, end time.Time) ([]Booking, error)
	GetBookingsByBookable(ctx context.Context, bookableID uuid.UUID, start, end time.Time) ([]Booking, error)
	CreateBooking(ctx context.Context, booking *Booking) error
	UpdateBookingStatus(ctx context.Context, bookingID uuid.UUID, status BookingStatus, updatedBy uuid.UUID) error
}

type AuditRepository interface {
	RecordLog(ctx context.Context, log *AuditLog) error
	GetLogsByEntity(ctx context.Context, entityType string, entityID uuid.UUID) ([]AuditLog, error)
}
