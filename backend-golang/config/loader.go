package config

import (
	"log"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	Port        string
	DBHost      string
	DBPort      string
	DBUser      string
	DBPassword  string
	DBName      string
	DBSSLMode   string
	Environment string
}

var AppConfig Config

func LoadConfig() {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	// If no .env file, just use environment variables
	if err := viper.ReadInConfig(); err != nil {
		log.Println("[CONFIG] No .env file found, using system environment variables")
	}

	AppConfig = Config{
		Port:        getEnv("PORT", "8080"),
		DBHost:      getEnv("DB_HOST", "localhost"),
		DBPort:      getEnv("DB_PORT", "3306"),
		DBUser:      getEnv("DB_USER", "root"),
		DBPassword:  getEnv("DB_PASSWORD", ""),
		DBName:      getEnv("DB_NAME", "iam_governance"),
		DBSSLMode:   getEnv("DB_SSL_MODE", "disable"),
		Environment: getEnv("ENVIRONMENT", "development"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	if value := viper.GetString(key); value != "" {
		return value
	}
	return defaultValue
}
