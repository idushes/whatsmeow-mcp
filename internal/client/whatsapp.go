package client

import (
	"time"
	"whatsmeow-mcp/internal/types"
)

// WhatsAppClient manages WhatsApp client state and operations
type WhatsAppClient struct {
	isLoggedIn bool
	qrCode     string
	messages   []types.Message
}

// NewWhatsAppClient creates a new WhatsApp client instance
func NewWhatsAppClient() *WhatsAppClient {
	return &WhatsAppClient{
		isLoggedIn: false,
		qrCode:     "2@ABC123DEF456GHI789JKL012MNO345PQR678STU901VWX234YZ,567890abcdef1234567890abcdef1234567890ab,cdefghijklmnopqrstuvwxyz",
		messages: []types.Message{
			{
				ID:        "msg_001",
				From:      "1234567890@s.whatsapp.net",
				Text:      "Hello! How are you?",
				Timestamp: time.Now().Unix() - 3600,
				Chat:      "1234567890@s.whatsapp.net",
			},
			{
				ID:        "msg_002",
				From:      "9876543210@s.whatsapp.net",
				Text:      "How are you doing?",
				Timestamp: time.Now().Unix() - 1800,
				Chat:      "9876543210@s.whatsapp.net",
			},
		},
	}
}

// IsLoggedIn returns the current authentication status
func (c *WhatsAppClient) IsLoggedIn() bool {
	return c.isLoggedIn
}

// GetQRCode returns the current QR code for authentication
func (c *WhatsAppClient) GetQRCode() string {
	return c.qrCode
}

// SetLoggedIn updates the authentication status
func (c *WhatsAppClient) SetLoggedIn(status bool) {
	c.isLoggedIn = status
}

// UpdateQRCode generates a new QR code
func (c *WhatsAppClient) UpdateQRCode() {
	c.qrCode = "2@XYZ789ABC012DEF345GHI678JKL901MNO234PQR567STU890VWX123YZ,456789abcdef0123456789abcdef0123456789ab,cdefghijklmnopqrstuvwxyz"
}

// AddMessage adds a new message to the client's message history
func (c *WhatsAppClient) AddMessage(message types.Message) {
	c.messages = append(c.messages, message)
}

// GetChatMessages returns messages for a specific chat
func (c *WhatsAppClient) GetChatMessages(chatJID string, count int, beforeMessageID string) []types.Message {
	var chatMessages []types.Message

	for _, msg := range c.messages {
		if msg.Chat == chatJID {
			if beforeMessageID == "" || msg.ID != beforeMessageID {
				chatMessages = append(chatMessages, msg)
			}
		}
	}

	// Limit results
	if count > 0 && len(chatMessages) > count {
		chatMessages = chatMessages[:count]
	}

	return chatMessages
}

// GetAllMessages returns all messages in the client
func (c *WhatsAppClient) GetAllMessages() []types.Message {
	return c.messages
}
