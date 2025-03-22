package main

import (
	"kingdup/db"
	"kingdup/routes"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func init() {
	// Only load .env file if running locally
	env := os.Getenv("ENV")
	if env == "" || env == "development" {
		if err := godotenv.Load(); err != nil {
			log.Println("⚠️ Could not load .env file (expected in local dev)")
		} else {
			log.Println("✅ .env loaded (local dev)")
		}
	}
}

func main() {
	// Init DB
	db.Init()

	// Init Router
	router := routes.RegisterRoutes()

	// Port
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 Server running on port %s\n", port)
	router.Run(":" + port)
}
