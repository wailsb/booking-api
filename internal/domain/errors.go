package domain

import "errors"

var (
	ErrInvalidTimeRange = errors.New("start time must be before end time")
	ErrBookingInPast    = errors.New("cannot book in the past")
	ErrDoubleBooking    = errors.New("time slot is already booked")
	ErrBookableNotFound = errors.New("bookable resource not found")
	ErrBookingNotFound  = errors.New("booking not found")
)
