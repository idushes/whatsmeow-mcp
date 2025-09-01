package client

import "whatsmeow-mcp/internal/types"

// WhatsAppClientInterface defines the interface that all WhatsApp clients must implement
type WhatsAppClientInterface interface {
	// Authentication methods
	IsLoggedIn() bool
	GetQRCode() string
	Connect() error

	// Message methods
	SendMessage(to, text, quotedMessageID string) (*types.MessageResponse, error)
	GetChatMessages(chatJID string, count int, beforeMessageID string) []types.Message
	GetAllMessages() []types.Message
	AddMessage(message types.Message)

	// Contact methods
	IsOnWhatsApp(phones []string) ([]types.WhatsAppCheckResult, error)

	// Compatibility methods (for backward compatibility with mock client)
	SetLoggedIn(status bool)
	UpdateQRCode()
}
