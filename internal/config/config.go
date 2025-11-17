package config

import (
	"fmt"
	"os"
)

// Config holds all application configuration
type Config struct {
	// GitHub configuration
	GitHubToken   string
	WebhookSecret string

	// Slack configuration
	SlackToken   string
	SlackChannel string

	// AI configuration
	GeminiAPIKey string  // CHANGED FROM AnthropicAPIKey

	// Application configuration
	Environment string
	Port        string
	LogLevel    string
}

func Load() (*Config, error) {
	cfg := &Config{
		GitHubToken:   os.Getenv("GITHUB_TOKEN"),
		WebhookSecret: os.Getenv("WEBHOOK_SECRET"),
		SlackToken:    os.Getenv("SLACK_TOKEN"),
		SlackChannel:  os.Getenv("SLACK_CHANNEL"),
		GeminiAPIKey:  os.Getenv("GEMINI_API_KEY"),  // CHANGED
		Environment:   getEnvOrDefault("ENVIRONMENT", "development"),
		Port:          getEnvOrDefault("PORT", "8080"),
		LogLevel:      getEnvOrDefault("LOG_LEVEL", "info"),
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) Validate() error {
	if c.GitHubToken == "" {
		return fmt.Errorf("GITHUB_TOKEN is required")
	}
	if c.WebhookSecret == "" {
		return fmt.Errorf("WEBHOOK_SECRET is required")
	}
	if c.SlackToken == "" {
		return fmt.Errorf("SLACK_TOKEN is required")
	}
	if c.SlackChannel == "" {
		return fmt.Errorf("SLACK_CHANNEL is required")
	}
	if c.GeminiAPIKey == "" {
		return fmt.Errorf("GEMINI_API_KEY is required")
	}
	return nil
}

// IsDevelopment returns true if running in development mode
func (c *Config) IsDevelopment() bool {
	return c.Environment == "development"
}

// IsProduction returns true if running in production mode
func (c *Config) IsProduction() bool {
	return c.Environment == "production"
}

// getEnvOrDefault gets an environment variable or returns a default value
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}