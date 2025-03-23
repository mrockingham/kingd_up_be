package db

import (
	"database/sql"
	"kingdup/models"
)

// CreateOrderItem inserts a single order item
func CreateOrderItem(db *sql.DB, item *models.OrderItem) error {
	query := `
		INSERT INTO order_items (
			order_id, variant_id, product_name, size, color, quantity, price_each
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := db.Exec(
		query,
		item.OrderID,
		item.VariantID,
		item.ProductName,
		item.Size,
		item.Color,
		item.Quantity,
		item.PriceEach,
	)

	return err
}
