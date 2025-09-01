-- Remove delivery status field from messages table
DROP INDEX IF EXISTS messages_delivery_status_idx;
ALTER TABLE messages DROP COLUMN IF EXISTS is_delivered;
