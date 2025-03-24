package auth

import (
	"kingdup/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func MeHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		var user models.User
		if err := db.First(&user, userID).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"id":          user.ID,
			"name":        user.Name,
			"email":       user.Email,
			"is_verified": user.IsVerified,
			"is_admin":    user.IsAdmin,
		})
	}
}
