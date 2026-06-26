package config

import (
	"log"
	"os"
	"strconv"

	"github.com/spf13/viper"
)

type Config struct {
	Port             string
	DBHost           string
	DBPort           string
	DBUser           string
	DBPassword       string
	DBName           string
	DBSSLMode        string
	Environment      string
	AppName          string
	LogToStdout      bool
	LogToFile        bool
	LogRetentionDays int
	LogLevel         string
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
		Port:             getEnv("PORT", "8080"),
		DBHost:           getEnv("DB_HOST", "localhost"),
		DBPort:           getEnv("DB_PORT", "3306"),
		DBUser:           getEnv("DB_USER", "root"),
		DBPassword:       getEnv("DB_PASSWORD", ""),
		DBName:           getEnv("DB_NAME", "iam_governance"),
		DBSSLMode:        getEnv("DB_SSL_MODE", "disable"),
		Environment:      getEnv("ENVIRONMENT", "development"),
		AppName:          getEnv("APP_NAME", "iam-governance"),
		LogToStdout:      getEnvBool("LOG_TO_STDOUT", true),
		LogToFile:        getEnvBool("LOG_TO_FILE", true),
		LogRetentionDays: getEnvInt("LOG_RETENTION_DAYS", 21),
		LogLevel:         getEnv("LOG_LEVEL", "info"),
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

func getEnvBool(key string, defaultValue bool) bool {
	valStr := getEnv(key, "")
	if valStr == "" {
		return defaultValue
	}
	val, err := strconv.ParseBool(valStr)
	if err != nil {
		return defaultValue
	}
	return val
}

func getEnvInt(key string, defaultValue int) int {
	valStr := getEnv(key, "")
	if valStr == "" {
		return defaultValue
	}
	val, err := strconv.Atoi(valStr)
	if err != nil {
		return defaultValue
	}
	return val
}

