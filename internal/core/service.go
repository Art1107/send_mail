package core

import (
    "sendgridtest/internal/domain"
    "sendgridtest/internal/ports"
    "sendgridtest/pkg/logger"
    "time"
)

type EventService struct {
    notifier ports.Notifier
    logger   *logger.Logger
}

func NewEventService(notifier ports.Notifier, logger *logger.Logger) *EventService {
    return &EventService{
        notifier: notifier,
        logger:   logger,
    }
}

const (
    EventDelivered  = "delivered"
    EventOpen       = "open"
    EventClick      = "click"
    EventBounce     = "bounce"
    EventSpamReport = "spam_report"
)

func isMainEvent(eventType string) bool {
    switch eventType {
    case EventDelivered, EventOpen, EventClick, EventBounce, EventSpamReport:
        return true
    default:
        return false
    }
}

func (s *EventService) HandleEvent(event domain.SendgridEvent) error {
    if !isMainEvent(event.Event) {
        return nil
    }

    s.logEvent(event)

    if event.Event == EventBounce || event.Event == EventSpamReport {
        return s.handleNegativeEvent(event)
    }

    return nil
}

func (s *EventService) logEvent(event domain.SendgridEvent) {
    s.logger.Info("SendGrid Event",
        "event", event.Event,
        "email", event.Email,
        "timestamp", time.Unix(event.Timestamp, 0).Format("2006-01-02 15:04:05"))
}

func (s *EventService) handleNegativeEvent(event domain.SendgridEvent) error {
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
    return nil
}