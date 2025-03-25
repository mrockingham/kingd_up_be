package main

import (
	"kingdup/db"
	"kingdup/routes"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate"
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

func runMigrations() {
	m, err := migrate.New(
		"file://migrations", // relative path to your migrations
		os.Getenv("DATABASE_URL"),
	)
	if err != nil {
		log.Printf("❌ Migration setup failed: %v\n", err)
		return
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Printf("❌ Migration failed: %v\n", err)
		return
	}

	log.Println("✅ Migrations ran successfully")
}

func main() {

	db.Init()

	runMigrations()

	r := gin.Default()
	// Init DB

	// Register routes
	router := routes.RegisterRoutes(r, db.DB)

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
