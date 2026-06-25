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

// DB connection string helper
func getDSN() string {
	user := getEnv("DB_USER", "root")
	pass := getEnv("DB_PASSWORD", "secret")
	host := getEnv("DB_HOST", "127.0.0.1")
	port := getEnv("DB_PORT", "3306")
	name := getEnv("DB_NAME", "folksart_iam")
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&multiStatements=true", user, pass, host, port, name)
}

// InitDB initializes the MySQL connection pool
func InitDB() *sql.DB {
	dsn := getDSN()
	host := getEnv("DB_HOST", "127.0.0.1")
	port := getEnv("DB_PORT", "3306")

	var db *sql.DB
	var err error

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

// End of file
