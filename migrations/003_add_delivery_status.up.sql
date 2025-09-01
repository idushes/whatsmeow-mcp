-- Add delivery status field to messages table
ALTER TABLE messages ADD COLUMN is_delivered BOOLEAN NOT NULL DEFAULT false;

-- Create index for delivered messages
CREATE INDEX messages_delivery_status_idx ON messages(is_delivered) WHERE is_delivered = true;

-- Add comments for clarity
COMMENT ON COLUMN messages.is_delivered IS 'Whether the message has been delivered to the recipient';
COMMENT ON COLUMN messages.is_read IS 'Whether the message has been read by the recipient';
