package service

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

// PaymentAccount represents a user's payment account
type PaymentAccount struct {
	ID       int
	UserID   int
	Username string
	Email    string
	Balance  float64
}

func CreatePaymentAccount(userID int, username, email string, db *sql.DB) error {
	// Check if the user already has a payment account
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM payment_accounts WHERE user_id=$1)", userID).Scan(&exists)
	if err != nil {
		log.Printf("Error checking existing account for user %d: %v", userID, err)
		return err
	}
	if exists {
		log.Printf("Payment account for user %d already exists, skipping creation", userID)
		return nil
	}

	// Insert new payment account
	query := `
        INSERT INTO payment_accounts (user_id, username, email, balance)
        VALUES ($1, $2, $3, $4)
        RETURNING id
    `
	var id int
	err = db.QueryRow(query, userID, username, email, 0.0).Scan(&id)
	if err != nil {
		log.Printf("Failed to create payment account for user %d: %v", userID, err)
		return err
	}

	log.Printf("Created payment account %d for user %d (%s, %s)", id, userID, username, email)
	return nil
}

func InitializeDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open DB: %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping DB: %w", err)
	}
	return db, nil
}
