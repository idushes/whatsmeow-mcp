package types

// SendMessageParams represents parameters for sending a text message
type SendMessageParams struct {
	To              string `json:"to" description:"WhatsApp JID of recipient. For phone numbers: 'phonenumber@s.whatsapp.net' (e.g. '1234567890@s.whatsapp.net'). For groups: 'groupid@g.us'"`
	Text            string `json:"text" description:"Text content of the message to send (plain text, no formatting)"`
	QuotedMessageID string `json:"quoted_message_id,omitempty" description:"Optional message ID to reply to. Use message ID from previous chat history to quote/reply to that message"`
}

// IsOnWhatsappParams represents parameters for checking WhatsApp registration status
type IsOnWhatsappParams struct {
	Phones []string `json:"phones" description:"Array of phone numbers in international format (e.g., +1234567890) to check"`
}

// GetChatHistoryParams represents parameters for retrieving chat message history
type GetChatHistoryParams struct {
	Chat            string `json:"chat" description:"WhatsApp JID (chat identifier) to retrieve messages from"`
	Count           int    `json:"count,omitempty" description:"Maximum number of messages to retrieve (default: 50, max: 100)"`
	BeforeMessageID string `json:"before_message_id,omitempty" description:"Optional message ID to retrieve messages before this point (for pagination)"`
}

// GetUnreadMessagesParams represents parameters for retrieving unread messages
type GetUnreadMessagesParams struct {
	Chat  string `json:"chat,omitempty" description:"Optional WhatsApp JID to filter unread messages from a specific chat. If omitted, returns unread messages from all chats"`
	Count int    `json:"count,omitempty" description:"Maximum number of unread messages to retrieve (default: 50, max: 100)"`
}

// MarkMessagesAsReadParams represents parameters for marking messages as read
type MarkMessagesAsReadParams struct {
	Chat string `json:"chat" description:"WhatsApp JID (chat identifier) to mark messages as read in this chat"`
}
