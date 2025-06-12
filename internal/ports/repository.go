package ports

import "sendgridtest/internal/domain"

type EventRepository interface {
	SaveEvent(event domain.SendgridEvent) error
	GetEventsByEmail(email string) ([]domain.SendgridEvent, error)
}
