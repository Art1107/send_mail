package adapters

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"send_mail/domain"
)

type LarkNotifier struct {
	webhookURL string
}

func NewLarkNotifier(webhookURL string) *LarkNotifier {
	return &LarkNotifier{webhookURL: webhookURL}
}

func (l *LarkNotifier) Notify(event domain.SendgridEvent) error {
	categoryStr := fmt.Sprintf("%v", event.Category)
	message := map[string]interface{}{
		"msg_type": "text",
		"content": map[string]string{
			"text": fmt.Sprintf("SendGrid:\nEvent: %s\nEmail: %s\nCategory: %s",
				event.Event, event.Email, categoryStr),
		},
	}

	body, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("error marshaling message: %v", err)
	}

	resp, err := http.Post(l.webhookURL, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("lark notify failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("lark response status: %s", resp.Status)
	}
	return nil
}
