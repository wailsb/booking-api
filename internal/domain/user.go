package domain

import (
	"time"

	"github.com/google/uuid"
)

type UserRole string

const (
	RoleAdmin   UserRole = "ADMIN"
	RoleManager UserRole = "MANAGER"
	RoleStaff   UserRole = "STAFF"
)

type User struct {
	ID           uuid.UUID `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"` // Omit from JSON serialization
	FullName     string    `json:"full_name"`
	Role         UserRole  `json:"role"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}