package domain

type SendgridEvent struct {
	Email            string `json:"email"`
	Timestamp        int64  `json:"timestamp"`
	Event            string `json:"event"`
	Category         any    `json:"category"`
	SGEventID        string `json:"sg_event_id"`
	SGMessageID      string `json:"sg_message_id"`
	Response         string `json:"response,omitempty"`
	Reason           string `json:"reason,omitempty"`
	Status           string `json:"status,omitempty"`
	Type             string `json:"type,omitempty"`
	IP               string `json:"ip,omitempty"`
	UserAgent        string `json:"useragent,omitempty"`
	URL              string `json:"url,omitempty"`
	ASMGroupID       int    `json:"asm_group_id,omitempty"`
	NewsletterID     int    `json:"newsletter_id,omitempty"`
	NewsletterSendID int    `json:"newsletter_send_id,omitempty"`
}
