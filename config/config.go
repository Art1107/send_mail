package config

import (
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort        string
	LarkWebhookURL    string
	LogFile           string
	SendgridPublicKey string
}

type LarkConfig struct {
	WebhookURL     string
	Timeout        time.Duration
	MaxRetries     int
	RateLimit      float64
	MaxMessageSize int
}

const (
	DefaultServerPort = ":8080"
	DefaultLogFile    = "sendgrid_events.log"
	DefaultWebhookURL = "https://open.larksuite.com/open-apis/bot/v2/hook/662a59e9-d3b5-41db-81c0-911e083525d5"
)

func init() {
	loadEnvFile()
}

func loadEnvFile() {
	envLocations := []string{".env", "../.env", "../../.env"}

	for _, loc := range envLocations {
		if err := godotenv.Load(filepath.Clean(loc)); err == nil {
			log.Printf("Loaded .env from: %s", loc)
			return
		}
	}

	log.Printf("Warning: Could not load .env file from any location")
}

func NewConfig() *Config {
	publicKey := os.Getenv("SENDGRID_PUBLIC_KEY")
	logPublicKeyInfo(publicKey)

	return &Config{
		ServerPort:        getEnvOrDefault("SERVER_PORT", DefaultServerPort),
		LarkWebhookURL:    getEnvOrDefault("LARK_WEBHOOK_URL", DefaultWebhookURL),
		LogFile:           getEnvOrDefault("LOG_FILE", DefaultLogFile),
		SendgridPublicKey: publicKey,
	}
}

func NewLarkConfig() LarkConfig {
	return LarkConfig{
		WebhookURL:     getEnvOrDefault("LARK_WEBHOOK_URL", ""),
		Timeout:        getDurationOrDefault("LARK_TIMEOUT", 5*time.Second),
		MaxRetries:     getIntOrDefault("LARK_MAX_RETRIES", 3),
		RateLimit:      getFloat64OrDefault("LARK_RATE_LIMIT", 10),
		MaxMessageSize: getIntOrDefault("LARK_MAX_MESSAGE_SIZE", 1000),
	}
}

func logPublicKeyInfo(key string) {
	log.Printf("Debug: SENDGRID_PUBLIC_KEY length: %d", len(key))
	if len(key) > 0 {
		maxLen := min(len(key), 20)
		log.Printf("Debug: SENDGRID_PUBLIC_KEY first chars: %s", key[:maxLen])
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func getDurationOrDefault(key string, defaultValue time.Duration) time.Duration {
	if value, err := time.ParseDuration(os.Getenv(key)); err == nil {
		return value
	}
	return defaultValue
}

func getIntOrDefault(key string, defaultValue int) int {
	if value, err := strconv.Atoi(os.Getenv(key)); err == nil {
		return value
	}
	return defaultValue
}

func getFloat64OrDefault(key string, defaultValue float64) float64 {
	if value, err := strconv.ParseFloat(os.Getenv(key), 64); err == nil {
		return value
	}
	return defaultValue
}
