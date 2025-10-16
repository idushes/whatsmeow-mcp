package client

import (
	"context"
	"whatsmeow-mcp/internal/types"
)

// WhatsAppClientInterface defines the interface that all WhatsApp clients must implement
type WhatsAppClientInterface interface {
	// Authentication methods
	IsLoggedIn() bool
	GetQRCode() string
	Connect() error

	// Message methods
	SendMessage(ctx context.Context, to, text, quotedMessageID string) (*types.MessageResponse, error)
	GetChatMessages(chatJID string, count int, beforeMessageID string) []types.Message
	GetUnreadMessages(chatJID string, count int) []types.Message
	GetAllMessages() []types.Message
	AddMessage(message types.Message)
	MarkMessagesAsRead(chatJID string) error

	// Contact methods
	IsOnWhatsApp(phones []string) ([]types.WhatsAppCheckResult, error)

	// Subscription methods
	GetSubscriptionManager() *SubscriptionManager
}
