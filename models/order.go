package models

import "time"

type Order struct {
	ID         int         `json:"id" gorm:"column:id"`
	UserID     *int        `json:"user_id,omitempty" gorm:"column:user_id"`
	GuestEmail *string     `json:"guest_email,omitempty" gorm:"column:guest_email"`
	GuestName  *string     `json:"guest_name,omitempty" gorm:"column:guest_name"`
	Status     string      `json:"status" gorm:"column:status"`
	Total      float64     `json:"total" gorm:"column:total_amount"`
	CreatedAt  time.Time   `json:"created_at" gorm:"column:created_at"`
	Items      []OrderItem `json:"items" gorm:"foreignKey:OrderID"`
}
