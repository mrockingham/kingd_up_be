package models

type OrderItem struct {
	ID          int     `json:"id" gorm:"column:id"`
	OrderID     int     `json:"order_id" gorm:"column:order_id"`
	VariantID   int64   `json:"variant_id" gorm:"column:variant_id"`
	ProductName string  `json:"product_name" gorm:"column:product_name"`
	Size        string  `json:"size" gorm:"column:size"`
	Color       string  `json:"color" gorm:"column:color"`
	Quantity    int     `json:"quantity" gorm:"column:quantity"`
	PriceEach   float64 `json:"price_each" gorm:"column:price_each"`
}
