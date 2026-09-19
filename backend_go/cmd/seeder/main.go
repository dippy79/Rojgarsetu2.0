package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL is required")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error opening DB: %v", err)
	}
	defer db.Close()

	// 1. Seed Roles/Static Data (if any)
	// 2. Seed Admin User
	email := "admin@rojgarsetu.in"
	password := "Admin@12345678" // Secure default
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	_, err = db.Exec(`
		INSERT INTO users (email, password_hash, role, name, is_active)
		VALUES ($1, $2, 'admin', 'System Administrator', true)
		ON CONFLICT (email) DO NOTHING`,
		email, string(hashedPassword))

	if err != nil {
		log.Fatalf("Error seeding admin: %v", err)
	}

	log.Println("Database seeding completed successfully.")
}
