package lark

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sendgridtest/internal/domain"
	"time"
)

type Notifier struct {
	webhookURL string
	client     *http.Client
	timeout    time.Duration
}

type larkMessage struct {
	MsgType string            `json:"msg_type"`
	Content map[string]string `json:"content"`
}

func NewNotifier(webhookURL string) *Notifier {
	return &Notifier{
		webhookURL: webhookURL,
		client: &http.Client{
			Timeout: 10 * time.Second,
			Transport: &http.Transport{
				MaxIdleConns:        100,
				MaxIdleConnsPerHost: 100,
				IdleConnTimeout:     90 * time.Second,
				TLSHandshakeTimeout: 10 * time.Second,
			},
		},
		timeout: 5 * time.Second,
	}
}

func (n *Notifier) Notify(event domain.SendgridEvent) error {
	ctx, cancel := context.WithTimeout(context.Background(), n.timeout)
	defer cancel()

	// ตรวจสอบความถูกต้องของข้อมูล
	if err := n.validateEvent(event); err != nil {
		return fmt.Errorf("invalid event data: %w", err)
	}

	message := larkMessage{
		MsgType: "text",
		Content: map[string]string{
			"text": n.sanitizeMessage(event),
		},
	}

	payload, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", n.webhookURL, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "SendGrid-Webhook-Handler/1.0")

	resp, err := n.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send notification: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return domain.ErrNotificationError
	}

	return nil
}

// validateEvent ตรวจสอบความถูกต้องของข้อมูล event
func (n *Notifier) validateEvent(event domain.SendgridEvent) error {
	if event.Email == "" {
		return fmt.Errorf("empty email address")
	}

	if event.Event == "" {
		return fmt.Errorf("empty event type")
	}

	// ตรวจสอบประเภท event ที่ยอมรับ
	validEvents := map[string]bool{
		"bounce":      true,
		"spam_report": true,
	}

	if !validEvents[event.Event] {
		return fmt.Errorf("unsupported event type: %s", event.Event)
	}

	return nil
}

// sanitizeMessage ทำความสะอาดข้อความก่อนส่ง
func (n *Notifier) sanitizeMessage(event domain.SendgridEvent) string {
	maxLength := 1000 // จำกัดความยาวข้อความ

	msg := fmt.Sprintf("Email Event: %s\nEmail: %s\nTimestamp: %d",
		event.Event,
		event.Email,
		event.Timestamp,
	)

	if len(msg) > maxLength {
		msg = msg[:maxLength] + "..."
	}

	return msg
}
