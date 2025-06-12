package lark

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sendgridtest/internal/domain"
)

type Notifier struct {
	webhookURL string
	client     *http.Client
}

type larkMessage struct {
	MsgType string            `json:"msg_type"`
	Content map[string]string `json:"content"`
}

func NewNotifier(webhookURL string) *Notifier {
	return &Notifier{
		webhookURL: webhookURL,
		client:     &http.Client{},
	}
}

func (n *Notifier) Notify(event domain.SendgridEvent) error {
	message := larkMessage{
		MsgType: "text",
		Content: map[string]string{
			"text": fmt.Sprintf("Email Event: %s\nEmail: %s", event.Event, event.Email),
		},
	}

	payload, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	resp, err := n.client.Post(n.webhookURL, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("failed to send notification: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("notification failed with status: %d", resp.StatusCode)
	}

	return nil
}
