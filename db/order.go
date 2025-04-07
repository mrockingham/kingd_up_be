package db

import (
	"context"
	"database/sql"
	"fmt"
	"kingdup/models"
	"time"
)

type OrderInput struct {
	UserID          *int64
	GuestName       *string
	GuestEmail      *string
	TotalAmount     float64
	Status          string
	ShippingAddress string
	Items           []OrderItemInput
}

type OrderItemInput struct {
	VariantID   int64
	ProductName string
	Size        string
	Color       string
	Quantity    int
	PriceEach   float64
}

// CreateOrder inserts a new order into the database
func CreateOrder(db *sql.DB, order *models.Order) (int64, error) {
	query := `
		INSERT INTO orders (user_id, guest_email, guest_name, status, total_price, created_at)
		VALUES ($1, $2, $3, $4, $5, $6) RETURNING id
	`

	var orderID int64
	err := db.QueryRow(
		query,
		order.UserID,
		order.GuestEmail,
		order.GuestName,
		order.Status,
		order.Total,
		time.Now(),
	).Scan(&orderID)

	return orderID, err
}

func CreateOrderWithItems(ctx context.Context, db *sql.DB, input OrderInput) (int64, error) {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	var orderID int64
	err = tx.QueryRowContext(ctx, `
	INSERT INTO orders (user_id, guest_name, guest_email, shipping_address, total_amount, status)
	VALUES ($1, $2, $3, $4, $5, $6) RETURNING id
`, input.UserID, input.GuestName, input.GuestEmail, input.ShippingAddress, input.TotalAmount, input.Status).Scan(&orderID)
	if err != nil {
		return 0, fmt.Errorf("failed to create order: %w", err)
	}

	for _, item := range input.Items {
		_, err := tx.ExecContext(ctx, `
			INSERT INTO order_items (order_id, variant_id, product_name, size, color, quantity, price_each)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
		`, orderID, item.VariantID, item.ProductName, item.Size, item.Color, item.Quantity, item.PriceEach)
		if err != nil {
			return 0, fmt.Errorf("failed to insert order item: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("commit failed: %w", err)
	}

	return orderID, nil
}

func SaveGuestOrderFromCheckout(ctx context.Context, db *sql.DB, items []OrderItemInput, guestEmail string) (int64, error) {
	total := 0.0
	for _, item := range items {
		total += item.PriceEach * float64(item.Quantity)
	}

	order := OrderInput{
		UserID:          nil,
		GuestEmail:      &guestEmail,
		Status:          "pending",
		TotalAmount:     total,
		ShippingAddress: "", // You can wire this in if available
		Items:           items,
	}

	return CreateOrderWithItems(ctx, db, order)
}
