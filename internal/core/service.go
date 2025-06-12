package core

import (
	"sendgridtest/internal/domain"
	"sendgridtest/internal/ports"
	"sendgridtest/pkg/logger"
	"strings"
	"time"
)

type EventService struct {
	repository ports.EventRepository
	notifier   ports.Notifier
	logger     *logger.Logger
}

func NewEventService(repo ports.EventRepository, notifier ports.Notifier, logger *logger.Logger) *EventService {
	return &EventService{
		repository: repo,
		notifier:   notifier,
		logger:     logger,
	}
}

func isMainEvent(eventType string) bool {
	switch eventType {
	case "delivered", "open", "click", "bounce", "spam_report":
		return true
	default:
		return false
	}
}

func (s *EventService) HandleEvent(event domain.SendgridEvent) error {

	if !isMainEvent(event.Event) {
		return nil
	}

	s.logger.Info("SendGrid Event",
		"event", event.Event,
		"email", event.Email,
		"timestamp", time.Unix(event.Timestamp, 0).Format("2006-01-02 15:04:05"))

	if err := s.repository.SaveEvent(event); err != nil {
		s.logger.Error("Failed to save event",
			"error", err,
			"event", event.Event,
			"email", event.Email)
		return err
	}

	// Notify only for critical events
	if event.Event == "bounce" || event.Event == "spam_report" {
		if err := s.notifier.Notify(event); err != nil {
			s.logger.Error("Failed to send notification",
				"error", err,
				"event", event.Event,
				"email", event.Email)
			return err
		}
		s.logger.Info("Notification sent",
			"event", event.Event,
			"email", event.Email)
	}

	return nil
}

func cleanEmailAddress(email string) string {

	if idx := strings.Index(email, "."); idx != -1 {
		email = email[:idx]
	}
	return email
}
