package types

// StandardResponse represents the common response structure for all tools
type StandardResponse struct {
	Success bool        `json:"success"`
	Error   *ErrorInfo  `json:"error,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// ErrorInfo represents detailed error information
type ErrorInfo struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// LoginStatusResponse represents the response for authentication status check
type LoginStatusResponse struct {
	LoggedIn bool `json:"logged_in"`
	Success  bool `json:"success"`
}

// QRCodeResponse represents the response for QR code generation
type QRCodeResponse struct {
	QRCode    string `json:"qr_code"`    // Raw QR code string
	Code      string `json:"code"`       // Same as qr_code (for compatibility)
	ImageURL  string `json:"image_url"`  // URL to QR code image
	Timeout   int    `json:"timeout"`    // Timeout in seconds
	Success   bool   `json:"success"`    // Success status
	ExpiresAt int64  `json:"expires_at"` // Unix timestamp when QR expires
}

// MessageResponse represents the response for message sending
type MessageResponse struct {
	MessageID       string `json:"message_id"`
	Timestamp       int64  `json:"timestamp"`
	Success         bool   `json:"success"`
	To              string `json:"to"`
	Text            string `json:"text"`
	QuotedMessageID string `json:"quoted_message_id,omitempty"`
}

// WhatsAppCheckResult represents a single phone number check result
type WhatsAppCheckResult struct {
	Phone        string `json:"phone"`
	IsOnWhatsApp bool   `json:"is_on_whatsapp"`
	JID          string `json:"jid"`
}

// WhatsAppCheckResponse represents the response for WhatsApp registration check
type WhatsAppCheckResponse struct {
	Results []WhatsAppCheckResult `json:"results"`
	Success bool                  `json:"success"`
}

// Message represents a single chat message
type Message struct {
	ID              string `json:"id"`
	From            string `json:"from"`
	To              string `json:"to,omitempty"`
	Text            string `json:"text"`
	Timestamp       int64  `json:"timestamp"`
	Chat            string `json:"chat"`
	QuotedMessageID string `json:"quoted_message_id,omitempty"`
}

// ChatHistoryResponse represents the response for chat history retrieval
type ChatHistoryResponse struct {
	Messages []Message `json:"messages"`
	HasMore  bool      `json:"has_more"`
	Success  bool      `json:"success"`
	Chat     string    `json:"chat"`
	Count    int       `json:"count"`
}
