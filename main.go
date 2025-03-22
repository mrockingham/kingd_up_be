package main

import (
	"kingdup/api"
	"log"
	"os"

	"kingdup/db"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	// Only load .env in development
	if os.Getenv("ENV") != "production" {
		err := godotenv.Load()
		if err != nil {
			log.Println("⚠️ Could not load .env file (expected in development only)")
		}
	}
}

func main() {
	// Load .env
	if err := godotenv.Load(); err != nil {
		log.Fatal("❌ Error loading .env file")
	}

	// Initialize DB
	db.Init()

	// Set up Gin
	router := gin.Default()

	// Enable CORS (add this)
	router.Use(cors.Default())

	// Test route
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})
	// Sync products
	router.GET("/api/sync-products", api.SyncProductsHandler)

	// Get products
	router.GET("/api/products", api.GetProductsHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 Server running on port %s\n", port)
	router.Run(":" + port)
}
