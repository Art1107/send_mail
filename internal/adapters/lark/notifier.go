package lark

import (
	"bytes"
	"encoding/json"
	"net/http"
	"sendgridtest/internal/domain"
)

type Notifier struct {
	webhookURL string
}

func NewNotifier(webhookURL string) *Notifier {
	return &Notifier{
		webhookURL: webhookURL,
	}
}

func (n *Notifier) Notify(event domain.SendgridEvent) error {
	message := map[string]interface{}{
		"msg_type": "text",
		"content": map[string]string{
			"text": "Email Event: " + event.Event + "\nEmail: " + event.Email,
		},
	}

	payload, err := json.Marshal(message)
	if err != nil {
		return err
	}

	_, err = http.Post(n.webhookURL, "application/json", bytes.NewBuffer(payload))
	return err
}
