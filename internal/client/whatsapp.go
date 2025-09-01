package client

import (
	"fmt"
	"time"
	"whatsmeow-mcp/internal/types"
)

// WhatsAppClient manages WhatsApp client state and operations (mock implementation)
type WhatsAppClient struct {
	isLoggedIn bool
	qrCode     string
	messages   []types.Message
}

// Ensure WhatsAppClient implements WhatsAppClientInterface
var _ WhatsAppClientInterface = (*WhatsAppClient)(nil)

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

// Connect establishes connection (mock implementation)
func (c *WhatsAppClient) Connect() error {
	return nil // Mock always succeeds
}

// SendMessage sends a message (mock implementation)
func (c *WhatsAppClient) SendMessage(to, text, quotedMessageID string) (*types.MessageResponse, error) {
	messageId := fmt.Sprintf("msg_%d", time.Now().Unix())
	timestamp := time.Now().Unix()

	response := &types.MessageResponse{
		MessageID:       messageId,
		Timestamp:       timestamp,
		Success:         true,
		To:              to,
		Text:            text,
		QuotedMessageID: quotedMessageID,
	}

	// Add message to history
	newMessage := types.Message{
		ID:              messageId,
		From:            "self",
		To:              to,
		Text:            text,
		Timestamp:       timestamp,
		Chat:            to,
		QuotedMessageID: quotedMessageID,
	}
	c.AddMessage(newMessage)

	return response, nil
}

// IsOnWhatsApp checks phone numbers (mock implementation)
func (c *WhatsAppClient) IsOnWhatsApp(phones []string) ([]types.WhatsAppCheckResult, error) {
	results := make([]types.WhatsAppCheckResult, 0, len(phones))
	for _, phone := range phones {
		// Check phone number registration status
		isRegistered := len(phone) > 10 && phone[len(phone)-1] != '0'

		results = append(results, types.WhatsAppCheckResult{
			Phone:        phone,
			IsOnWhatsApp: isRegistered,
			JID:          phone + "@s.whatsapp.net",
		})
	}
	return results, nil
}
