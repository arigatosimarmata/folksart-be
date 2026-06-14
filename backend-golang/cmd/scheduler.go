package cmd

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/robfig/cron/v3"
	"github.com/spf13/cobra"
)

var schedulerCmd = &cobra.Command{
	Use:   "scheduler",
	Short: "Start the background task scheduler",
	Run: func(cmd *cobra.Command, args []string) {
		runScheduler()
	},
}

func init() {
	rootCmd.AddCommand(schedulerCmd)
}

func runScheduler() {
	log.Println("[SCHEDULER] Initializing background task jobs...")

	c := cron.New()

	// Example job: Run every minute
	_, err := c.AddFunc("* * * * *", func() {
		log.Println("[JOB] Heartbeat: Scheduler is active and processing tasks...")
	})
	if err != nil {
		log.Fatalf("[SCHEDULER] Failed to schedule heartbeat job: %v", err)
	}

	// Example: Run at midnight for maintenance
	_, err = c.AddFunc("0 0 * * *", func() {
		log.Println("[JOB] Daily Maintenance: Cleaning up expired sessions and audit logs...")
		// Add cleanup logic here
	})

	c.Start()
	log.Println("[SCHEDULER] Cron engine started successfully.")

	// Wait for termination signal
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	<-sig

	log.Println("[SCHEDULER] Shutting down cron engine...")
	c.Stop()
}
