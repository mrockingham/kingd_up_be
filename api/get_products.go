package api

import (
	"kingdup/db"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetProductsHandler(c *gin.Context) {
	var products []db.Product

	// Preload the associated variants
	if err := db.DB.Preload("Variants").Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch products"})
		return
	}

	c.JSON(http.StatusOK, products)
}
