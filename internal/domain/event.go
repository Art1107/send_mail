package domain

type SendgridEvent struct {
	Email       string   `json:"email"`
	Timestamp   int64    `json:"timestamp"`
	Event       string   `json:"event"`
	Category    []string `json:"category,omitempty"`
	SGEventID   string   `json:"sg_event_id,omitempty"`
	SGMessageID string   `json:"sg_message_id,omitempty"`
	IP          string   `json:"ip,omitempty"`
	UserAgent   string   `json:"useragent,omitempty"`
	Status      string   `json:"status,omitempty"`
}
