package api

import (
	"kingdup/db"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GET /api/products
func GetProductsHandler(c *gin.Context) {
	var products []db.Product
	err := db.DB.Preload("Variants").Find(&products).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load products"})
		return
	}
	c.JSON(http.StatusOK, products)
}

// GET /api/products/:id
func GetProductHandler(c *gin.Context) {
	idParam := c.Param("id")
	productID, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	var product db.Product
	err = db.DB.Preload("Variants").First(&product, "printful_id = ?", productID).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	c.JSON(http.StatusOK, product)
}
