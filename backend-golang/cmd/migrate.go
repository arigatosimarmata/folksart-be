package cmd

import (
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/spf13/cobra"
	"react-example/backend-golang/config"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Run database migrations",
	Run: func(cmd *cobra.Command, args []string) {
		runMigrations()
	},
}

func init() {
	rootCmd.AddCommand(migrateCmd)
}

func runMigrations() {
	cfg := config.AppConfig
	dsn := fmt.Sprintf("mysql://%s:%s@tcp(%s:%s)/%s?parseTime=true", 
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)
	
	m, err := migrate.New("file://backend-golang/migrations", dsn)
	if err != nil {
		log.Fatalf("[MIGRATE] Could not initialize migration: %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("[MIGRATE] Migration failed: %v", err)
	}

	log.Println("[MIGRATE] Database migrations applied successfully!")
}
