package models

import "time"

type User struct {
	ID           int       `json:"id" gorm:"column:id"`
	Email        string    `json:"email" gorm:"column:email"`
	Name         string    `json:"name" gorm:"column:name"`
	PasswordHash string    `json:"-" gorm:"column:password_hash"`
	Phone        string    `json:"phone" gorm:"column:phone"`
	IsVerified   bool      `json:"is_verified" gorm:"column:is_verified"`
	IsAdmin      bool      `json:"is_admin" gorm:"column:is_admin"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}
