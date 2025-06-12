package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"send_mail/adapters"
	"send_mail/config"
	"send_mail/domain"
)

type App struct {
	config   *config.Config
	notifier *adapters.LarkNotifier
}

func NewApp() *App {
	cfg := config.NewConfig()
	return &App{
		config:   cfg,
		notifier: adapters.NewLarkNotifier(cfg.LarkWebhookURL),
	}
}

func (app *App) handleWebhook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		log.Printf("❌ Invalid method: %s", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("❌ Failed to read body: %v", err)
		http.Error(w, "Cannot read body", http.StatusInternalServerError)
		return
	}

	var events []domain.SendgridEvent
	if err := json.Unmarshal(bodyBytes, &events); err != nil {
		log.Printf("❌ Invalid JSON: %v", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	for _, event := range events {
		log.Printf("📨 Event: %s | Email: %s | Status: %s | Time: %d",
			event.Event,
			event.Email,
			event.Status,
			event.Timestamp,
		)

		if event.Event == "spamreport" || event.Event == "dropped" || event.Event == "bounce" {
			log.Printf("⚠️ Important event detected: %s", event.Event)
			if err := app.notifier.Notify(event); err != nil {
				log.Printf("❌ Failed to send notification: %v", err)
			} else {
				log.Printf("✅ Notification sent successfully")
			}
		}
	}

	w.WriteHeader(http.StatusOK)
}

func main() {
	f, err := os.OpenFile("sendgrid_events.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)

	app := NewApp()
	http.HandleFunc("/webhook", app.handleWebhook)

	log.Println("🚀 Server started at", app.config.ServerPort)
	if err := http.ListenAndServe(app.config.ServerPort, nil); err != nil {
		log.Fatal(err)
	}
}
