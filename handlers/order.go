package handlers

import (
	"encoding/json"
	"kingdup/db"
	"kingdup/models"
	"net/http"

	"gorm.io/gorm"
)

type CreateOrderRequest struct {
	UserID     *int               `json:"user_id,omitempty"` // nullable
	GuestName  string             `json:"guest_name,omitempty"`
	GuestEmail string             `json:"guest_email,omitempty"`
	Items      []models.OrderItem `json:"items"`
	Total      float64            `json:"total"`
}

func toPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func toOrderItemInputs(items []models.OrderItem) []db.OrderItemInput {
	var result []db.OrderItemInput
	for _, item := range items {
		result = append(result, db.OrderItemInput{
			VariantID:   int64(item.VariantID),
			ProductName: item.ProductName,
			Size:        item.Size,
			Color:       item.Color,
			Quantity:    item.Quantity,
			PriceEach:   item.PriceEach,
		})
	}
	return result
}

func CreateOrderHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreateOrderRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		// Create order record
		order := models.Order{
			UserID:     req.UserID,
			GuestName:  toPtr(req.GuestName),
			GuestEmail: toPtr(req.GuestEmail),
			Status:     "pending",
			Total:      req.Total,
		}

		if err := db.Create(&order).Error; err != nil {
			http.Error(w, "Failed to create order", http.StatusInternalServerError)
			return
		}

		// Create order items
		for _, item := range req.Items {
			orderItem := models.OrderItem{
				OrderID:     order.ID,
				VariantID:   int64(item.VariantID),
				ProductName: item.ProductName,
				Size:        item.Size,
				Color:       item.Color,
				Quantity:    item.Quantity,
				PriceEach:   item.PriceEach,
			}
			if err := db.Create(&orderItem).Error; err != nil {
				http.Error(w, "Failed to create order item", http.StatusInternalServerError)
				return
			}
		}

		resp := map[string]interface{}{
			"message":  "Order created successfully",
			"order_id": order.ID,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}
