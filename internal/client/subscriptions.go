package client

import (
	"context"
	"sync"

	"github.com/mark3labs/mcp-go/server"
)

// SubscriptionManager manages chat subscriptions per MCP session
type SubscriptionManager struct {
	// sessionID -> chatJID -> subscribed
	subscriptions map[string]map[string]bool
	mutex         sync.RWMutex
	mcpServer     *server.MCPServer
}

// NewSubscriptionManager creates a new subscription manager
func NewSubscriptionManager(mcpServer *server.MCPServer) *SubscriptionManager {
	return &SubscriptionManager{
		subscriptions: make(map[string]map[string]bool),
		mcpServer:     mcpServer,
	}
}

// Subscribe adds a chat subscription for a session
func (sm *SubscriptionManager) Subscribe(sessionID, chatJID string) bool {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// Initialize session map if it doesn't exist
	if sm.subscriptions[sessionID] == nil {
		sm.subscriptions[sessionID] = make(map[string]bool)
	}

	// Check if already subscribed
	if sm.subscriptions[sessionID][chatJID] {
		return false // Already subscribed
	}

	// Add subscription
	sm.subscriptions[sessionID][chatJID] = true
	return true // New subscription
}

// Unsubscribe removes a chat subscription for a session
func (sm *SubscriptionManager) Unsubscribe(sessionID, chatJID string) bool {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	if sm.subscriptions[sessionID] == nil {
		return false
	}

	if !sm.subscriptions[sessionID][chatJID] {
		return false // Not subscribed
	}

	delete(sm.subscriptions[sessionID], chatJID)

	// Clean up empty session map
	if len(sm.subscriptions[sessionID]) == 0 {
		delete(sm.subscriptions, sessionID)
	}

	return true
}

// IsSubscribed checks if a session is subscribed to a chat
func (sm *SubscriptionManager) IsSubscribed(sessionID, chatJID string) bool {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()

	if sm.subscriptions[sessionID] == nil {
		return false
	}

	return sm.subscriptions[sessionID][chatJID]
}

// GetSubscribedSessions returns all sessions subscribed to a specific chat
func (sm *SubscriptionManager) GetSubscribedSessions(chatJID string) []string {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()

	sessions := []string{}
	for sessionID, chats := range sm.subscriptions {
		if chats[chatJID] {
			sessions = append(sessions, sessionID)
		}
	}

	return sessions
}

// GetSubscriptions returns all chat subscriptions for a session
func (sm *SubscriptionManager) GetSubscriptions(sessionID string) []string {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()

	if sm.subscriptions[sessionID] == nil {
		return []string{}
	}

	chats := []string{}
	for chatJID := range sm.subscriptions[sessionID] {
		chats = append(chats, chatJID)
	}

	return chats
}

// CleanupSession removes all subscriptions for a session
func (sm *SubscriptionManager) CleanupSession(sessionID string) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	delete(sm.subscriptions, sessionID)
}

// NotifyNewMessage sends notification to all subscribed sessions about a new message
func (sm *SubscriptionManager) NotifyNewMessage(chatJID, messageID, from, text string, timestamp int64) {
	subscribedSessions := sm.GetSubscribedSessions(chatJID)

	if len(subscribedSessions) == 0 {
		return
	}

	notification := map[string]any{
		"method": "notifications/message",
		"params": map[string]any{
			"chat":       chatJID,
			"message_id": messageID,
			"from":       from,
			"text":       text,
			"timestamp":  timestamp,
		},
	}

	// Send notification to each subscribed session
	for _, sessionID := range subscribedSessions {
		if sm.mcpServer != nil {
			ctx := context.Background()
			err := sm.mcpServer.SendNotificationToClient(ctx, sessionID, notification)
			if err != nil {
				// Log error but continue with other sessions
				continue
			}
		}
	}
}

// NotifyMessageStatus sends notification about message delivery/read status
func (sm *SubscriptionManager) NotifyMessageStatus(chatJID, messageID, status string) {
	subscribedSessions := sm.GetSubscribedSessions(chatJID)

	if len(subscribedSessions) == 0 {
		return
	}

	notification := map[string]any{
		"method": "notifications/message",
		"params": map[string]any{
			"chat":       chatJID,
			"message_id": messageID,
			"status":     status, // "delivered" or "read"
		},
	}

	for _, sessionID := range subscribedSessions {
		if sm.mcpServer != nil {
			ctx := context.Background()
			sm.mcpServer.SendNotificationToClient(ctx, sessionID, notification)
		}
	}
}
