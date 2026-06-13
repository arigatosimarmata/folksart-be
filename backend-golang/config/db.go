package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

// InitDB initializes the MySQL connection pool with resilient connection-ready loops
func InitDB() *sql.DB {
	user := getEnv("DB_USER", "root")
	pass := getEnv("DB_PASSWORD", "secret")
	host := getEnv("DB_HOST", "127.0.0.1")
	port := getEnv("DB_PORT", "3306")
	name := getEnv("DB_NAME", "folksart_iam")

	// Create MySQL Connection String
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&multiStatements=true", user, pass, host, port, name)

	var db *sql.DB
	var err error

	log.Printf("[IAM-DB] Attempting connection on database host %s:%s...", host, port)

	// Resilient connection pattern (retry on cold starts e.g. in docker-compose)
	for i := 1; i <= 5; i++ {
		db, err = sql.Open("mysql", dsn)
		if err == nil {
			err = db.Ping()
			if err == nil {
				log.Println("[IAM-DB] MySQL Database connection verified and active!")
				break
			}
		}
		log.Printf("[IAM-DB] Attempt %d/5 failed. Retrying in 4 seconds... (%v)", i, err)
		time.Sleep(4 * time.Second)
	}

	if err != nil {
		log.Fatalf("[IAM-DB] CRITICAL: Could not establish database pool: %v", err)
	}

	// Set optimal connection pool properties
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(10 * time.Minute)

	DB = db
	return db
}

// Helpers
func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
