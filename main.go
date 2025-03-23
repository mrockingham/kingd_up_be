package main

import (
	"kingdup/db"
	"kingdup/routes"
	"log"
	"os"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
)

func init() {
	// Load .env in local/dev environment
	env := os.Getenv("ENV")
	if env == "" || env == "development" {
		if err := godotenv.Load(); err != nil {
			log.Println("⚠️ Could not load .env file (expected in local dev)")
		} else {
			log.Println("✅ .env loaded")
		}
	}
}

func main() {

	// Init DB
	db.Init()

	// Register routes
	router := routes.RegisterRoutes()

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 Server running on port %s\n", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}
}
