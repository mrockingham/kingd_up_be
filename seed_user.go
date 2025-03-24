// seed_user.go
package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

func seedUser() {
	connStr := "postgres://postgres:password@localhost:5432/kingdup_db?sslmode=disable" // update this
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("❌ Failed to connect to DB:", err)
	}
	defer db.Close()

	email := "test@example.com"
	password := "password123"
	name := "Test User"

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal("❌ Failed to hash password:", err)
	}

	_, err = db.Exec(`
		INSERT INTO users (email, name, password_hash, phone, is_verified, is_admin)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, email, name, string(hashedPassword), "", true, false)
	if err != nil {
		log.Fatal("❌ Failed to insert test user:", err)
	}

	fmt.Println("✅ Test user created! Email:", email, "Password:", password)
}
