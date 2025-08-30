package types

// SendMessageParams represents parameters for sending a text message
type SendMessageParams struct {
	To              string `json:"to" description:"WhatsApp JID (phone number with @s.whatsapp.net suffix) of the recipient"`
	Text            string `json:"text" description:"Text content of the message to send"`
	QuotedMessageID string `json:"quoted_message_id,omitempty" description:"Optional ID of a previous message to quote/reply to"`
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

// DummyParams represents a placeholder parameter for tools that don't need parameters
type DummyParams struct {
	RandomString string `json:"random_string" description:"Dummy parameter for no-parameter tools"`
}
