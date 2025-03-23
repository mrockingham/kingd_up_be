package routes

import (
	"kingdup/api"
	order "kingdup/handlers/order" // Use an alias

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(r *gin.Engine, db *gorm.DB) *gin.Engine {

	router := r

	// Enable CORS
	router.Use(cors.Default())

	// Test
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	// Sync products
	router.GET("/api/sync-products", api.SyncProductsHandler)

	// Get all products
	router.GET("/api/products", api.GetProductsHandler)

	// Get product by ID
	router.GET("/api/products/:id", api.GetProductHandler)

	orderGroup := router.Group("/orders")
	{
		orderGroup.POST("/", order.CreateOrderHandler(db))
		// Add more as needed
	}

	return router
}
