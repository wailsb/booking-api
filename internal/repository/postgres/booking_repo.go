package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"booking-api/internal/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BookingRepository struct {
	db *pgxpool.Pool
}

func NewBookingRepository(db *pgxpool.Pool) domain.BookingRepository {
	return &BookingRepository{db: db}
}

func (r *BookingRepository) CreateBookable(ctx context.Context, b *domain.Bookable) error {
	metadataJSON, err := json.Marshal(b.Metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	query := `
		INSERT INTO bookables (id, name, description, type, metadata, created_by, updated_by, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	_, err = r.db.Exec(ctx, query,
		b.ID, b.Name, b.Description, b.Type, metadataJSON,
		b.CreatedBy, b.UpdatedBy, b.CreatedAt, b.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create bookable: %w", err)
	}

	return nil
}

func (r *BookingRepository) GetBookableByID(ctx context.Context, id uuid.UUID) (*domain.Bookable, error) {
	query := `
		SELECT id, name, description, type, metadata, created_by, updated_by, created_at, updated_at
		FROM bookables WHERE id = $1
	`

	var b domain.Bookable
	var metadataRaw []byte

	err := r.db.QueryRow(ctx, query, id).Scan(
		&b.ID, &b.Name, &b.Description, &b.Type, &metadataRaw,
		&b.CreatedBy, &b.UpdatedBy, &b.CreatedAt, &b.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrBookableNotFound
		}
		return nil, fmt.Errorf("failed to get bookable: %w", err)
	}

	if len(metadataRaw) > 0 {
		_ = json.Unmarshal(metadataRaw, &b.Metadata)
	}

	return &b, nil
}

func (r *BookingRepository) CreateBooking(ctx context.Context, b *domain.Booking) error {
	timeRangeStr := fmt.Sprintf("[%s, %s)",
		b.StartTime.Format(time.RFC3339Nano),
		b.EndTime.Format(time.RFC3339Nano),
	)

	query := `
		INSERT INTO bookings (
			id, bookable_id, customer_name, customer_email, customer_phone,
			booking_window, status, created_by, updated_by, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6::tstzrange, $7, $8, $9, $10, $11)
	`

	_, err := r.db.Exec(ctx, query,
		b.ID, b.BookableID, b.CustomerName, b.CustomerEmail, b.CustomerPhone,
		timeRangeStr, b.Status, b.CreatedBy, b.UpdatedBy, b.CreatedAt, b.UpdatedAt,
	)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23P01" {
			return domain.ErrDoubleBooking
		}
		return fmt.Errorf("failed to create booking: %w", err)
	}

	return nil
}

func (r *BookingRepository) GetBookingByID(ctx context.Context, id uuid.UUID) (*domain.Booking, error) {
	query := `
		SELECT id, bookable_id, customer_name, customer_email, customer_phone,
		       lower(booking_window) AS start_time, upper(booking_window) AS end_time,
		       status, created_by, updated_by, created_at, updated_at
		FROM bookings WHERE id = $1
	`

	var b domain.Booking
	err := r.db.QueryRow(ctx, query, id).Scan(
		&b.ID, &b.BookableID, &b.CustomerName, &b.CustomerEmail, &b.CustomerPhone,
		&b.StartTime, &b.EndTime, &b.Status, &b.CreatedBy, &b.UpdatedBy,
		&b.CreatedAt, &b.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrBookingNotFound
		}
		return nil, fmt.Errorf("failed to get booking: %w", err)
	}

	return &b, nil
}

func (r *BookingRepository) ListBookingsByBookable(ctx context.Context, bookableID uuid.UUID, start, end time.Time) ([]domain.Booking, error) {
	timeRangeStr := fmt.Sprintf("[%s, %s)", start.Format(time.RFC3339Nano), end.Format(time.RFC3339Nano))

	query := `
		SELECT id, bookable_id, customer_name, customer_email, customer_phone,
		       lower(booking_window) AS start_time, upper(booking_window) AS end_time,
		       status, created_by, updated_by, created_at, updated_at
		FROM bookings
		WHERE bookable_id = $1 AND booking_window && $2::tstzrange AND status != 'CANCELLED'
		ORDER BY booking_window ASC
	`

	rows, err := r.db.Query(ctx, query, bookableID, timeRangeStr)
	if err != nil {
		return nil, fmt.Errorf("failed to list bookings: %w", err)
	}
	defer rows.Close()

	var bookings []domain.Booking
	for rows.Next() {
		var b domain.Booking
		err := rows.Scan(
			&b.ID, &b.BookableID, &b.CustomerName, &b.CustomerEmail, &b.CustomerPhone,
			&b.StartTime, &b.EndTime, &b.Status, &b.CreatedBy, &b.UpdatedBy,
			&b.CreatedAt, &b.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan booking: %w", err)
		}
		bookings = append(bookings, b)
	}

	return bookings, nil
}

// Satisfies duplicate interface method requirements
func (r *BookingRepository) GetBookingsByBookable(ctx context.Context, bookableID uuid.UUID, start, end time.Time) ([]domain.Booking, error) {
	return r.ListBookingsByBookable(ctx, bookableID, start, end)
}

func (r *BookingRepository) UpdateBookingStatus(ctx context.Context, bookingID uuid.UUID, status domain.BookingStatus, updatedBy uuid.UUID) error {
	query := `
		UPDATE bookings
		SET status = $1, updated_by = $2, updated_at = NOW()
		WHERE id = $3
	`

	res, err := r.db.Exec(ctx, query, status, updatedBy, bookingID)
	if err != nil {
		return fmt.Errorf("failed to update booking status: %w", err)
	}

	if res.RowsAffected() == 0 {
		return domain.ErrBookingNotFound
	}

	return nil
}
