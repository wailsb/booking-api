package middleware

import (
	"context"
	"net/http"
	"strings"
	"booking-api/internal/domain"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type contextKey string

const (
	UserIDKey contextKey = "user_id"
	RoleKey   contextKey = "role"
)

type Claims struct {
	UserID uuid.UUID `json:"user_id"`
	Role   domain.UserRole `json:"role"`
	jwt.RegisteredClaims
}

func AuthMiddleware(jwtSecret []byte, requiredRoles ...domain.UserRole) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				http.Error(w, `{"error": "Unauthorized: missing or invalid token format"}`, http.StatusUnauthorized)
				return
			}

			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
			claims := &Claims{}

			token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
				return jwtSecret, nil
			})

			if err != nil || !token.Valid {
				http.Error(w, `{"error": "Unauthorized: invalid or expired token"}`, http.StatusUnauthorized)
				return
			}

			// Role check (Y must have permission for Z)
			if len(requiredRoles) > 0 {
				hasRole := false
				for _, role := range requiredRoles {
					if claims.Role == role {
						hasRole = true
						break
					}
				}
				if !hasRole {
					http.Error(w, `{"error": "Forbidden: insufficient permissions"}`, http.StatusForbidden)
					return
				}
			}

			// Pass User ID (X) and Role (Y) into Context
			// Ensure UserIDKey uses the contextKey type explicitly in auth.go:
			ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
			ctx = context.WithValue(ctx, RoleKey, claims.Role) // Use 'ctx', not 'r.Context()'

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
