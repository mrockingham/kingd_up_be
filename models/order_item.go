package models

type OrderItem struct {
	ID          int     `json:"id" db:"id"`
	OrderID     int     `json:"order_id" db:"order_id"`
	VariantID   int64   `json:"variant_id" db:"variant_id"`
	ProductName string  `json:"product_name" db:"product_name"`
	Size        string  `json:"size" db:"size"`
	Color       string  `json:"color" db:"color"`
	Quantity    int     `json:"quantity" db:"quantity"`
	PriceEach   float64 `json:"price_each" db:"price_each"`
	TotalPrice  float64 `json:"total_price" db:"total_price"`
}
