package client

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"whatsmeow-mcp/internal/database"
	"whatsmeow-mcp/internal/types"

	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/store/sqlstore"
	waTypes "go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	"google.golang.org/protobuf/proto"
)

// WhatsmeowClient implements WhatsApp functionality using the whatsmeow library
type WhatsmeowClient struct {
	client    *whatsmeow.Client
	container *sqlstore.Container
	db        *sql.DB

	// Database store for messages
	messageStore *database.MessageStore

	// QR channel for receiving QR codes
	qrChan chan string

	// Current QR code
	currentQR string

	// Connection status
	connected bool
	loggedIn  bool

	// Our JID for database operations
	ourJID string
}

// Ensure WhatsmeowClient implements WhatsAppClientInterface
var _ WhatsAppClientInterface = (*WhatsmeowClient)(nil)

// NewWhatsmeowClient creates a new WhatsApp client using whatsmeow
func NewWhatsmeowClient(db *sql.DB) (*WhatsmeowClient, error) {
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

	wc := &WhatsmeowClient{
		client:       client,
		container:    container,
		db:           db,
		messageStore: database.NewMessageStore(db),
		qrChan:       make(chan string, 1),
		connected:    false,
		loggedIn:     false,
	}

	// Set our JID if device is already paired
	if deviceStore.ID != nil {
		wc.ourJID = deviceStore.ID.String()
		log.Printf("Device already paired. Our JID: %s", wc.ourJID)
	}

	// Set up event handlers
	wc.setupEventHandlers()

	// Try to connect automatically if already paired
	if deviceStore.ID != nil {
		log.Printf("Attempting to connect with existing session...")
		go func() {
			time.Sleep(1 * time.Second) // Give time for event handlers to be set up
			if err := wc.Connect(); err != nil {
				log.Printf("Failed to auto-connect: %v", err)
			}
		}()
	}

	return wc, nil
}

// setupEventHandlers configures event handlers for WhatsApp events
func (wc *WhatsmeowClient) setupEventHandlers() {
	wc.client.AddEventHandler(func(evt interface{}) {
		switch v := evt.(type) {
		case *events.Message:
			wc.handleMessage(v)
		case *events.QR:
			wc.handleQR(v)
		case *events.Connected:
			wc.connected = true
			log.Printf("Connected to WhatsApp")
			// If we have a stored session and are connected, we're logged in
			if wc.client.Store.ID != nil {
				wc.loggedIn = true
				if wc.ourJID == "" {
					wc.ourJID = wc.client.Store.ID.String()
				}
				log.Printf("Restored session. Logged in as: %s", wc.ourJID)
				// Request history sync for restored sessions
				go wc.requestHistorySync()
			}
		case *events.Disconnected:
			wc.connected = false
			log.Printf("Disconnected from WhatsApp")
		case *events.LoggedOut:
			wc.loggedIn = false
			wc.ourJID = ""
			log.Printf("Logged out from WhatsApp")
		case *events.PairSuccess:
			wc.loggedIn = true
			wc.ourJID = wc.client.Store.ID.String()
			log.Printf("Successfully paired with WhatsApp. Our JID: %s", wc.ourJID)
			// Request history sync after successful pairing
			go wc.requestHistorySync()
		case *events.HistorySync:
			wc.handleHistorySync(v)
		default:
			// Log other events for debugging
			log.Printf("Received event: %T", v)
		}
	})
}

// handleMessage processes incoming messages
func (wc *WhatsmeowClient) handleMessage(evt *events.Message) {
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

	// Save message to database
	if wc.ourJID != "" {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := wc.messageStore.SaveMessage(ctx, message, wc.ourJID); err != nil {
			log.Printf("Failed to save message to database: %v", err)
		}
	}

	log.Printf("Received message from %s: %s", message.From, message.Text)
}

// handleQR processes QR code events
func (wc *WhatsmeowClient) handleQR(evt *events.QR) {
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
func (wc *WhatsmeowClient) Connect() error {
	if wc.connected {
		return nil // Already connected
	}
	return wc.client.Connect()
}

// Disconnect closes the connection to WhatsApp
func (wc *WhatsmeowClient) Disconnect() {
	wc.client.Disconnect()
}

// IsLoggedIn returns the current authentication status
func (wc *WhatsmeowClient) IsLoggedIn() bool {
	return wc.client.IsConnected() && wc.client.IsLoggedIn()
}

// GetQRCode returns the current QR code for authentication
func (wc *WhatsmeowClient) GetQRCode() string {
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
func (wc *WhatsmeowClient) SendMessage(to, text, quotedMessageID string) (*types.MessageResponse, error) {
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

	// Add to message history in database
	message := types.Message{
		ID:              resp.ID,
		From:            "self",
		To:              to,
		Text:            text,
		Timestamp:       resp.Timestamp.Unix(),
		Chat:            to,
		QuotedMessageID: quotedMessageID,
	}

	// Save sent message to database
	if wc.ourJID != "" {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := wc.messageStore.SaveMessage(ctx, message, wc.ourJID); err != nil {
			log.Printf("Failed to save sent message to database: %v", err)
		}
	}

	return response, nil
}

// IsOnWhatsApp checks if phone numbers are registered on WhatsApp
func (wc *WhatsmeowClient) IsOnWhatsApp(phones []string) ([]types.WhatsAppCheckResult, error) {
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

// GetChatMessages returns messages for a specific chat from database
func (wc *WhatsmeowClient) GetChatMessages(chatJID string, count int, beforeMessageID string) []types.Message {
	if wc.ourJID == "" {
		return []types.Message{}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	chatMessages, err := wc.messageStore.GetChatMessages(ctx, wc.ourJID, chatJID, count, beforeMessageID)
	if err != nil {
		log.Printf("Failed to get chat messages from database: %v", err)
		return []types.Message{}
	}

	return chatMessages
}

// GetUnreadMessages returns unread messages from database
func (wc *WhatsmeowClient) GetUnreadMessages(chatJID string, count int) []types.Message {
	if wc.ourJID == "" {
		return []types.Message{}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	unreadMessages, err := wc.messageStore.GetUnreadMessages(ctx, wc.ourJID, chatJID, count)
	if err != nil {
		log.Printf("Failed to get unread messages from database: %v", err)
		return []types.Message{}
	}

	return unreadMessages
}

// GetAllMessages returns all messages from database
func (wc *WhatsmeowClient) GetAllMessages() []types.Message {
	if wc.ourJID == "" {
		return []types.Message{}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	allMessages, err := wc.messageStore.GetAllMessages(ctx, wc.ourJID)
	if err != nil {
		log.Printf("Failed to get all messages from database: %v", err)
		return []types.Message{}
	}

	return allMessages
}

// AddMessage adds a new message to the database
func (wc *WhatsmeowClient) AddMessage(message types.Message) {
	if wc.ourJID == "" {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := wc.messageStore.SaveMessage(ctx, message, wc.ourJID); err != nil {
		log.Printf("Failed to add message to database: %v", err)
	}
}

// requestHistorySync requests message history from WhatsApp servers
func (wc *WhatsmeowClient) requestHistorySync() {
	if !wc.IsLoggedIn() {
		log.Printf("Cannot request history sync: not logged in")
		return
	}

	// Wait a bit for the connection to stabilize
	time.Sleep(3 * time.Second)

	log.Printf("Requesting message history sync...")

	// Send presence to trigger history sync
	err := wc.client.SendPresence(waTypes.PresenceAvailable)
	if err != nil {
		log.Printf("Failed to send presence for history sync: %v", err)
		return
	}

	// Additionally, try to fetch some recent conversations
	// This can help trigger history sync in some cases
	time.Sleep(2 * time.Second)

	log.Printf("History sync request completed. Waiting for HistorySync events...")
}

// handleHistorySync processes history sync events from WhatsApp
func (wc *WhatsmeowClient) handleHistorySync(evt *events.HistorySync) {
	if wc.ourJID == "" {
		log.Printf("Cannot process history sync: ourJID not set")
		return
	}

	log.Printf("Received history sync with %d conversations", len(evt.Data.GetConversations()))

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	messageCount := 0

	// Process each conversation in the history sync
	for _, conversation := range evt.Data.GetConversations() {
		chatJID := conversation.GetId()
		log.Printf("Processing conversation: %s", chatJID)

		// Process messages in this conversation
		for _, historyMsg := range conversation.GetMessages() {
			if historyMsg.GetMessage() == nil {
				continue
			}

			webMsg := historyMsg.GetMessage()
			if webMsg.GetKey() == nil {
				continue
			}

			// Convert WhatsApp message to our internal format
			message := wc.convertHistoryMessage(webMsg, chatJID)
			if message.ID == "" {
				continue
			}

			// Save to database
			if err := wc.messageStore.SaveMessage(ctx, message, wc.ourJID); err != nil {
				log.Printf("Failed to save history message: %v", err)
				continue
			}

			messageCount++
		}
	}

	log.Printf("Successfully processed %d messages from history sync", messageCount)
}

// convertHistoryMessage converts a WhatsApp history message to our internal format
func (wc *WhatsmeowClient) convertHistoryMessage(webMsg *waProto.WebMessageInfo, chatJID string) types.Message {
	if webMsg.GetKey() == nil {
		return types.Message{}
	}

	key := webMsg.GetKey()
	message := types.Message{
		ID:        key.GetId(),
		Chat:      chatJID,
		Timestamp: int64(webMsg.GetMessageTimestamp()),
	}

	// Determine sender
	if key.GetFromMe() {
		message.From = "self"
		message.To = key.GetRemoteJid()
	} else {
		message.From = key.GetRemoteJid()
		if key.GetParticipant() != "" {
			message.From = key.GetParticipant()
		}
	}

	// Extract message text
	if webMsg.GetMessage() != nil {
		msg := webMsg.GetMessage()
		if msg.GetConversation() != "" {
			message.Text = msg.GetConversation()
		} else if msg.GetExtendedTextMessage() != nil {
			message.Text = msg.GetExtendedTextMessage().GetText()
			if msg.GetExtendedTextMessage().GetContextInfo() != nil &&
				msg.GetExtendedTextMessage().GetContextInfo().GetStanzaId() != "" {
				message.QuotedMessageID = msg.GetExtendedTextMessage().GetContextInfo().GetStanzaId()
			}
		}
	}

	return message
}

// TestAddMessage is a helper function to add test messages for debugging
func (wc *WhatsmeowClient) TestAddMessage() {
	if wc.ourJID == "" {
		log.Printf("Cannot add test message: ourJID not set")
		return
	}

	testMessage := types.Message{
		ID:        fmt.Sprintf("test_%d", time.Now().Unix()),
		From:      "test@s.whatsapp.net",
		Chat:      "test@s.whatsapp.net",
		Text:      "Test message from history sync",
		Timestamp: time.Now().Unix(),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := wc.messageStore.SaveMessage(ctx, testMessage, wc.ourJID); err != nil {
		log.Printf("Failed to save test message: %v", err)
	} else {
		log.Printf("Successfully saved test message: %s", testMessage.ID)
	}
}
