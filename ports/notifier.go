package ports

import "send_mail/domain"

type EventNotifier interface {
	Notify(event domain.SendgridEvent) error
}
