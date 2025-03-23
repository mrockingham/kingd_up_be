// handlers/order/create.go (or wherever your handler is)
package order

import (
	"kingdup/db"
	"kingdup/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CreateOrderRequest struct {
	UserID          *int               `json:"user_id,omitempty"`
	GuestName       string             `json:"guest_name,omitempty"`
	GuestEmail      string             `json:"guest_email,omitempty"`
	Items           []models.OrderItem `json:"items"`
	ShippingAddress string             `json:"shipping_address"`
	Total           float64            `json:"total"`
}

func toPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func toInt64Ptr(i *int) *int64 {
	if i == nil {
		return nil
	}
	v := int64(*i)
	return &v
}

func CreateOrderHandler(dbConn *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req CreateOrderRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
			return
		}

		var items []db.OrderItemInput
		for _, item := range req.Items {
			items = append(items, db.OrderItemInput{
				VariantID:   int64(item.VariantID),
				ProductName: item.ProductName,
				Size:        item.Size,
				Color:       item.Color,
				Quantity:    item.Quantity,
				PriceEach:   item.PriceEach,
			})
		}

		input := db.OrderInput{
			UserID:      toInt64Ptr(req.UserID),
			GuestName:   toPtr(req.GuestName),
			GuestEmail:  toPtr(req.GuestEmail),
			TotalAmount: req.Total,
			Status:      "pending",
			Items:       items,
		}

		sqlDB, err := dbConn.DB()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get database connection"})
			return
		}

		orderID, err := db.CreateOrderWithItems(c.Request.Context(), sqlDB, input)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to create order",
				"details": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message":  "Order created successfully",
			"order_id": orderID,
		})
	}
}
