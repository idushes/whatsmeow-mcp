-- Create messages table for storing WhatsApp message history
CREATE TABLE messages (
    id TEXT PRIMARY KEY,
    our_jid TEXT NOT NULL,
    chat_jid TEXT NOT NULL,
    sender_jid TEXT NOT NULL,
    recipient_jid TEXT,
    message_text TEXT NOT NULL DEFAULT '',
    timestamp BIGINT NOT NULL,
    message_type TEXT NOT NULL DEFAULT 'text',
    quoted_message_id TEXT,
    is_from_me BOOLEAN NOT NULL DEFAULT false,
    is_read BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create indexes for efficient querying
CREATE INDEX messages_our_jid_idx ON messages(our_jid);
CREATE INDEX messages_chat_jid_idx ON messages(chat_jid);
CREATE INDEX messages_timestamp_idx ON messages(timestamp DESC);
CREATE INDEX messages_chat_timestamp_idx ON messages(chat_jid, timestamp DESC);
CREATE INDEX messages_sender_jid_idx ON messages(sender_jid);
CREATE INDEX messages_quoted_message_id_idx ON messages(quoted_message_id) WHERE quoted_message_id IS NOT NULL;

-- Create partial index for unread messages
CREATE INDEX messages_unread_idx ON messages(chat_jid, timestamp DESC) WHERE is_read = false;

-- Add trigger to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_messages_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER messages_updated_at_trigger
    BEFORE UPDATE ON messages
    FOR EACH ROW
    EXECUTE FUNCTION update_messages_updated_at();
