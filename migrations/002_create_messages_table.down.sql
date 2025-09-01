-- Drop trigger and function
DROP TRIGGER IF EXISTS messages_updated_at_trigger ON messages;
DROP FUNCTION IF EXISTS update_messages_updated_at();

-- Drop indexes
DROP INDEX IF EXISTS messages_unread_idx;
DROP INDEX IF EXISTS messages_quoted_message_id_idx;
DROP INDEX IF EXISTS messages_sender_jid_idx;
DROP INDEX IF EXISTS messages_chat_timestamp_idx;
DROP INDEX IF EXISTS messages_timestamp_idx;
DROP INDEX IF EXISTS messages_chat_jid_idx;
DROP INDEX IF EXISTS messages_our_jid_idx;

-- Drop messages table
DROP TABLE IF EXISTS messages CASCADE;
