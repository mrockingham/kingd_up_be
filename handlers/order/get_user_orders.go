package order

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetUserOrdersHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDVal, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
			return
		}

		userID := userIDVal.(int)

		var orders []struct {
			ID     int
			Status string
			Total  float64
		}

		if err := db.Raw("SELECT id, status, total_amount as total FROM orders WHERE user_id = ?", userID).Scan(&orders).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch orders"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"orders": orders})
	}
}
