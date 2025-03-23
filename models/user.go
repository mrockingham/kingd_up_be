package models

import "time"

type User struct {
	ID        int       `json:"id" db:"id"`
	Email     string    `json:"email" db:"email"`
	Name      string    `json:"name" db:"name"`
	Password  string    `json:"password,omitempty" db:"password"` // Optional: omit in responses
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
