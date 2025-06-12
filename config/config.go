package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort        string
	DatabaseDSN       string
	LarkWebhookURL    string
	LogFile           string
	SendgridPublicKey string
}

var defaultConfig = map[string]string{
	"SERVER_PORT":      ":8080",
	"DATABASE_DSN":     "root@tcp(localhost:3306)/sendgrid_events?parseTime=true",
	"LARK_WEBHOOK_URL": "https://open.larksuite.com/open-apis/bot/v2/hook/0cb0f8a4-9c7f-41f9-b30f-db6f22ae47c8",
	"LOG_FILE":         "sendgrid_events.log",
}

func init() {
	locations := []string{
		".env",
		"../.env",
		"../../.env",
	}

	var loaded bool
	for _, loc := range locations {
		if err := godotenv.Load(loc); err == nil {
			loaded = true
			log.Printf("Loaded .env from: %s", loc)
			break
		}
	}

	if !loaded {
		log.Printf("Warning: Could not load .env file from any location")
	}
}

func NewConfig() *Config {
	publicKey := os.Getenv("SENDGRID_PUBLIC_KEY")
	log.Printf("Debug: SENDGRID_PUBLIC_KEY length: %d", len(publicKey))
	log.Printf("Debug: SENDGRID_PUBLIC_KEY first chars: %s", publicKey[:min(len(publicKey), 20)])

	return &Config{
		ServerPort:        getEnvOrDefault("SERVER_PORT", defaultConfig["SERVER_PORT"]),
		DatabaseDSN:       getEnvOrDefault("DATABASE_DSN", defaultConfig["DATABASE_DSN"]),
		LarkWebhookURL:    getEnvOrDefault("LARK_WEBHOOK_URL", defaultConfig["LARK_WEBHOOK_URL"]),
		LogFile:           getEnvOrDefault("LOG_FILE", "sendgrid_events.log"),
		SendgridPublicKey: publicKey,
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
