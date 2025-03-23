package models

import "time"

type Order struct {
	ID         int       `json:"id" db:"id"`
	UserID     *int      `json:"user_id,omitempty" db:"user_id"` // null for guest
	GuestEmail *string   `json:"guest_email,omitempty" db:"guest_email"`
	GuestName  *string   `json:"guest_name,omitempty" db:"guest_name"`
	Status     string    `json:"status" db:"status"`
	Total      float64   `json:"total" db:"total"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
}
