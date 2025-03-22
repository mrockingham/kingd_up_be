package routes

import (
	"kingdup/api"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes() *gin.Engine {
	router := gin.Default()

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
	router.GET("/api/products/:id", api.GetProductsHandler)

	return router
}
