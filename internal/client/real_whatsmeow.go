package client

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"whatsmeow-mcp/internal/types"

	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/store/sqlstore"
	waTypes "go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	"google.golang.org/protobuf/proto"
)

// RealWhatsmeowClient implements WhatsApp functionality using the real whatsmeow library
type RealWhatsmeowClient struct {
	client    *whatsmeow.Client
	container *sqlstore.Container
	db        *sql.DB

	// QR channel for receiving QR codes
	qrChan chan string

	// Current QR code
	currentQR string

	// Message history (in-memory cache)
	messages []types.Message

	// Connection status
	connected bool
	loggedIn  bool
}

// Ensure RealWhatsmeowClient implements WhatsAppClientInterface
var _ WhatsAppClientInterface = (*RealWhatsmeowClient)(nil)

// NewRealWhatsmeowClient creates a new WhatsApp client using real whatsmeow
func NewRealWhatsmeowClient(db *sql.DB) (*RealWhatsmeowClient, error) {
	// Create SQL store container with auto-upgrade
	container := sqlstore.NewWithDB(db, "postgres", nil)

	// Upgrade database schema to latest version
	err := container.Upgrade(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to upgrade database schema: %w", err)
	}

	// Get the first device or create a new one
	deviceStore, err := container.GetFirstDevice(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to get device store: %w", err)
	}

	if deviceStore == nil {
		deviceStore = container.NewDevice()
	}

	// Create WhatsApp client
	client := whatsmeow.NewClient(deviceStore, nil)

	wc := &RealWhatsmeowClient{
		client:    client,
		container: container,
		db:        db,
		qrChan:    make(chan string, 1),
		messages:  make([]types.Message, 0),
		connected: false,
		loggedIn:  false,
	}

	// Set up event handlers
	wc.setupEventHandlers()

	return wc, nil
}

// setupEventHandlers configures event handlers for WhatsApp events
func (wc *RealWhatsmeowClient) setupEventHandlers() {
	wc.client.AddEventHandler(func(evt interface{}) {
		switch v := evt.(type) {
		case *events.Message:
			wc.handleMessage(v)
		case *events.QR:
			wc.handleQR(v)
		case *events.Connected:
			wc.connected = true
			log.Printf("Connected to WhatsApp")
		case *events.Disconnected:
			wc.connected = false
			log.Printf("Disconnected from WhatsApp")
		case *events.LoggedOut:
			wc.loggedIn = false
			log.Printf("Logged out from WhatsApp")
		case *events.PairSuccess:
			wc.loggedIn = true
			log.Printf("Successfully paired with WhatsApp")
		}
	})
}

// handleMessage processes incoming messages
func (wc *RealWhatsmeowClient) handleMessage(evt *events.Message) {
	message := types.Message{
		ID:        evt.Info.ID,
		From:      evt.Info.Sender.String(),
		Chat:      evt.Info.Chat.String(),
		Timestamp: evt.Info.Timestamp.Unix(),
	}

	if evt.Message.GetConversation() != "" {
		message.Text = evt.Message.GetConversation()
	} else if evt.Message.GetExtendedTextMessage() != nil {
		message.Text = evt.Message.GetExtendedTextMessage().GetText()
		if evt.Message.GetExtendedTextMessage().GetContextInfo() != nil &&
			evt.Message.GetExtendedTextMessage().GetContextInfo().GetStanzaID() != "" {
			message.QuotedMessageID = evt.Message.GetExtendedTextMessage().GetContextInfo().GetStanzaID()
		}
	}

	wc.messages = append(wc.messages, message)
	log.Printf("Received message from %s: %s", message.From, message.Text)
}

// handleQR processes QR code events
func (wc *RealWhatsmeowClient) handleQR(evt *events.QR) {
	wc.currentQR = evt.Codes[0]
	select {
	case wc.qrChan <- evt.Codes[0]:
	default:
		// Channel is full, replace the QR code
		<-wc.qrChan
		wc.qrChan <- evt.Codes[0]
	}
	log.Printf("QR code received, scan it with your phone")
}

// Connect establishes connection to WhatsApp
func (wc *RealWhatsmeowClient) Connect() error {
	if wc.connected {
		return nil // Already connected
	}
	return wc.client.Connect()
}

// Disconnect closes the connection to WhatsApp
func (wc *RealWhatsmeowClient) Disconnect() {
	wc.client.Disconnect()
}

// IsLoggedIn returns the current authentication status
func (wc *RealWhatsmeowClient) IsLoggedIn() bool {
	return wc.client.IsConnected() && wc.client.IsLoggedIn()
}

// GetQRCode returns the current QR code for authentication
func (wc *RealWhatsmeowClient) GetQRCode() string {
	if !wc.client.IsConnected() {
		// Try to connect if not connected
		if err := wc.Connect(); err != nil {
			log.Printf("Failed to connect: %v", err)
			return ""
		}
	}

	if wc.client.IsLoggedIn() {
		return "" // Already logged in
	}

	// Wait for QR code with timeout
	select {
	case qr := <-wc.qrChan:
		return qr
	case <-time.After(5 * time.Second):
		return wc.currentQR
	}
}

// SendMessage sends a text message to the specified recipient
func (wc *RealWhatsmeowClient) SendMessage(to, text, quotedMessageID string) (*types.MessageResponse, error) {
	if !wc.IsLoggedIn() {
		return nil, fmt.Errorf("not logged in")
	}

	// Parse recipient JID
	jid, err := waTypes.ParseJID(to)
	if err != nil {
		return nil, fmt.Errorf("invalid recipient JID: %w", err)
	}

	// Create message
	msg := &waProto.Message{
		Conversation: proto.String(text),
	}

	// Add quoted message if specified
	if quotedMessageID != "" {
		msg.ExtendedTextMessage = &waProto.ExtendedTextMessage{
			Text: proto.String(text),
			ContextInfo: &waProto.ContextInfo{
				StanzaID: proto.String(quotedMessageID),
			},
		}
		msg.Conversation = nil
	}

	// Send message
	resp, err := wc.client.SendMessage(context.Background(), jid, msg)
	if err != nil {
		return nil, fmt.Errorf("failed to send message: %w", err)
	}

	// Create response
	response := &types.MessageResponse{
		MessageID:       resp.ID,
		Timestamp:       resp.Timestamp.Unix(),
		Success:         true,
		To:              to,
		Text:            text,
		QuotedMessageID: quotedMessageID,
	}

	// Add to message history
	message := types.Message{
		ID:              resp.ID,
		From:            "self",
		To:              to,
		Text:            text,
		Timestamp:       resp.Timestamp.Unix(),
		Chat:            to,
		QuotedMessageID: quotedMessageID,
	}
	wc.messages = append(wc.messages, message)

	return response, nil
}

// IsOnWhatsApp checks if phone numbers are registered on WhatsApp
func (wc *RealWhatsmeowClient) IsOnWhatsApp(phones []string) ([]types.WhatsAppCheckResult, error) {
	if !wc.IsLoggedIn() {
		return nil, fmt.Errorf("not logged in")
	}

	// Convert phone numbers to clean format
	cleanPhones := make([]string, len(phones))
	for i, phone := range phones {
		// Clean phone number
		cleanPhone := strings.ReplaceAll(phone, "+", "")
		cleanPhone = strings.ReplaceAll(cleanPhone, " ", "")
		cleanPhone = strings.ReplaceAll(cleanPhone, "-", "")
		cleanPhones[i] = cleanPhone
	}

	// Check registration status
	results, err := wc.client.IsOnWhatsApp(cleanPhones)
	if err != nil {
		return nil, fmt.Errorf("failed to check WhatsApp registration: %w", err)
	}

	// Convert results
	checkResults := make([]types.WhatsAppCheckResult, len(phones))
	for i, phone := range phones {
		checkResults[i] = types.WhatsAppCheckResult{
			Phone:        phone,
			IsOnWhatsApp: results[i].IsIn,
			JID:          results[i].JID.String(),
		}
	}

	return checkResults, nil
}

// GetChatMessages returns messages for a specific chat
func (wc *RealWhatsmeowClient) GetChatMessages(chatJID string, count int, beforeMessageID string) []types.Message {
	var chatMessages []types.Message

	for _, msg := range wc.messages {
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
func (wc *RealWhatsmeowClient) GetAllMessages() []types.Message {
	return wc.messages
}

// SetLoggedIn updates the authentication status (for compatibility)
func (wc *RealWhatsmeowClient) SetLoggedIn(status bool) {
	// This method is kept for compatibility but doesn't do anything
	// as the login status is managed by whatsmeow
}

// UpdateQRCode generates a new QR code (for compatibility)
func (wc *RealWhatsmeowClient) UpdateQRCode() {
	// This method is kept for compatibility but doesn't do anything
	// as QR codes are managed by whatsmeow
}

// AddMessage adds a new message to the client's message history (for compatibility)
func (wc *RealWhatsmeowClient) AddMessage(message types.Message) {
	wc.messages = append(wc.messages, message)
}
