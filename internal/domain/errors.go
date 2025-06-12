package domain

type Error string

func (e Error) Error() string { return string(e) }

const (
	ErrInvalidEvent      = Error("invalid event")
	ErrDatabaseError     = Error("database error")
	ErrNotificationError = Error("notification error")
)
