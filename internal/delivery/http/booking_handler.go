package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"booking-api/internal/delivery/http/middleware"
	"booking-api/internal/domain"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type BookingHandler struct {
	bookingUC domain.BookingUseCase // Direct domain interface
}

func NewBookingHandler(bookingUC domain.BookingUseCase) *BookingHandler {
	return &BookingHandler{
		bookingUC: bookingUC,
	}
}

// Request DTOs
type CreateBookingRequest struct {
	BookableID    string `json:"bookable_id"`
	CustomerName  string `json:"customer_name"`
	CustomerEmail string `json:"customer_email"`
	CustomerPhone string `json:"customer_phone"`
	StartTime     string `json:"start_time"`
	EndTime       string `json:"end_time"`
}

// ----------------------------------------------------------------------------
// 1. CreateBooking
// ----------------------------------------------------------------------------
func (h *BookingHandler) CreateBooking(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Extract context actor ID
	actorID, ok := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		writeJSONError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req CreateBookingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}

	// Conversions & Validations
	bookableUUID, err := uuid.Parse(req.BookableID)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid bookable_id UUID format")
		return
	}

	start, err := time.Parse(time.RFC3339, req.StartTime)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid start_time format (expected RFC3339)")
		return
	}

	end, err := time.Parse(time.RFC3339, req.EndTime)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid end_time format (expected RFC3339)")
		return
	}

	now := time.Now().UTC()
	booking := &domain.Booking{
		ID:            uuid.New(),
		BookableID:    bookableUUID,
		CustomerName:  req.CustomerName,
		CustomerEmail: req.CustomerEmail,
		CustomerPhone: req.CustomerPhone,
		StartTime:     start,
		EndTime:       end,
		Status:        domain.StatusConfirmed,
		CreatedBy:     &actorID,
		UpdatedBy:     &actorID,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	created, err := h.bookingUC.CreateBooking(r.Context(), booking)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrDoubleBooking):
			writeJSONError(w, http.StatusConflict, "Requested time slot is already booked")
		case errors.Is(err, domain.ErrBookableNotFound):
			writeJSONError(w, http.StatusNotFound, "Bookable resource not found")
		default:
			writeJSONError(w, http.StatusInternalServerError, "Failed to create booking")
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(created)
}

// ----------------------------------------------------------------------------
// 2. CancelBooking
// ----------------------------------------------------------------------------
func (h *BookingHandler) CancelBooking(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idStr := chi.URLParam(r, "id")
	bookingID, err := uuid.Parse(idStr)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid booking UUID format")
		return
	}

	if err := h.bookingUC.CancelBooking(r.Context(), bookingID); err != nil {
		if errors.Is(err, domain.ErrBookingNotFound) {
			writeJSONError(w, http.StatusNotFound, "Booking not found")
			return
		}
		writeJSONError(w, http.StatusInternalServerError, "Failed to cancel booking")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ----------------------------------------------------------------------------
// 3. GetBookingByID
// ----------------------------------------------------------------------------
func (h *BookingHandler) GetBookingByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idStr := chi.URLParam(r, "id")
	bookingID, err := uuid.Parse(idStr)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid booking UUID format")
		return
	}

	booking, err := h.bookingUC.GetBookingByID(r.Context(), bookingID)
	if err != nil {
		if errors.Is(err, domain.ErrBookingNotFound) {
			writeJSONError(w, http.StatusNotFound, "Booking not found")
			return
		}
		writeJSONError(w, http.StatusInternalServerError, "Failed to retrieve booking")
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(booking)
}

// ----------------------------------------------------------------------------
// 4. ListBookingsByBookable
// ----------------------------------------------------------------------------
func (h *BookingHandler) ListBookingsByBookable(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	bookableIDStr := chi.URLParam(r, "bookableID")
	bookableID, err := uuid.Parse(bookableIDStr)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid bookable UUID format")
		return
	}

	startStr := r.URL.Query().Get("start")
	endStr := r.URL.Query().Get("end")

	start, err := time.Parse(time.RFC3339, startStr)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid start query parameter")
		return
	}

	end, err := time.Parse(time.RFC3339, endStr)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid end query parameter")
		return
	}

	bookings, err := h.bookingUC.ListBookingsByBookable(r.Context(), bookableID, start, end)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Failed to list bookings")
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(bookings)
}

// Helper function
func writeJSONError(w http.ResponseWriter, statusCode int, message string) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}
