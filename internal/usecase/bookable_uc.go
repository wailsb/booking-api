package usecase

import (
	"context"

	"booking-api/internal/domain"

	"github.com/google/uuid"
)

type BookableUseCase interface {
	CreateBookable(ctx context.Context, bookable *domain.Bookable) (*domain.Bookable, error)
	GetBookableByID(ctx context.Context, id uuid.UUID) (*domain.Bookable, error)
	ListBookables(ctx context.Context, limit, offset int) ([]domain.Bookable, error)
	UpdateBookable(ctx context.Context, bookable *domain.Bookable) (*domain.Bookable, error)
}
