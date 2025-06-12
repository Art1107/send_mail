package ports

import "sendgridtest/internal/domain"

type Notifier interface {
	Notify(event domain.SendgridEvent) error
}
