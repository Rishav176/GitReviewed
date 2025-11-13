package main

import (
	"log"
	"net/http"

	//"os"

	"github.com/Rishav176/GitReviewed/internal/config"
	"github.com/Rishav176/GitReviewed/internal/handlers"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file in development
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	log.Printf("Starting GitReviewed in %s mode", cfg.Environment)

	// Create webhook handler
	handler := handlers.NewWebhookHandler(cfg)

	// Register routes
	http.HandleFunc("/webhook", handler.HandleWebhook)
	http.HandleFunc("/health", handler.HealthCheck)
	http.HandleFunc("/test-slack", handler.TestSlack)

	// Start server
	addr := ":" + cfg.Port
	log.Printf("Server listening on %s", addr)
	log.Printf("Webhook endpoint: http://localhost%s/webhook", addr)
	log.Printf("Health check: http://localhost%s/health", addr)
	log.Printf("Test Slack: http://localhost%s/test-slack", addr)

	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}