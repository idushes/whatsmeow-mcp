package database

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"strings"
	"time"

	"whatsmeow-mcp/internal/types"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

// Connect establishes a connection to PostgreSQL database using Database URL
func Connect(databaseURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	return db, nil
}

// RunMigrations runs database migrations
func RunMigrations(db *sql.DB, migrationsPath string) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create postgres driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", migrationsPath),
		"postgres",
		driver,
	)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}

// CreateDatabase creates the database if it doesn't exist using Database URL
func CreateDatabase(databaseURL string) error {
	// Parse the database URL
	parsedURL, err := url.Parse(databaseURL)
	if err != nil {
		return fmt.Errorf("failed to parse database URL: %w", err)
	}

	// Extract database name from path
	dbName := strings.TrimPrefix(parsedURL.Path, "/")
	if dbName == "" {
		return fmt.Errorf("database name not found in URL")
	}

	// Create connection URL to postgres database (without specific database)
	postgresURL := *parsedURL
	postgresURL.Path = "/postgres"

	db, err := sql.Open("postgres", postgresURL.String())
	if err != nil {
		return fmt.Errorf("failed to connect to postgres: %w", err)
	}
	defer db.Close()

	// Check if database exists
	var exists bool
	err = db.QueryRow("SELECT EXISTS(SELECT datname FROM pg_catalog.pg_database WHERE datname = $1)", dbName).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check if database exists: %w", err)
	}

	if !exists {
		// Create database
		_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s", dbName))
		if err != nil {
			return fmt.Errorf("failed to create database: %w", err)
		}
	}

	return nil
}

// HealthCheck checks if the database connection is healthy
func HealthCheck(db *sql.DB) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return fmt.Errorf("database health check failed: %w", err)
	}

	return nil
}

// MessageStore handles database operations for messages
type MessageStore struct {
	db *sql.DB
}

// NewMessageStore creates a new MessageStore instance
func NewMessageStore(db *sql.DB) *MessageStore {
	return &MessageStore{db: db}
}

// GetDB returns the underlying database connection for health checks
func (ms *MessageStore) GetDB() *sql.DB {
	return ms.db
}

// SaveMessage saves a message to the database
func (ms *MessageStore) SaveMessage(ctx context.Context, msg types.Message, ourJID string) error {
	query := `
		INSERT INTO messages (
			id, our_jid, chat_jid, sender_jid, recipient_jid, 
			message_text, timestamp, message_type, quoted_message_id, 
			is_from_me, is_read
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		ON CONFLICT (id) DO UPDATE SET
			message_text = EXCLUDED.message_text,
			is_read = EXCLUDED.is_read,
			updated_at = NOW()
	`

	// Определяем, прочитано ли сообщение:
	// - Исходящие сообщения (от нас) считаются прочитанными
	// - Входящие сообщения (от других) изначально непрочитанные
	isFromMe := msg.From == "self"
	isRead := isFromMe // Только наши сообщения считаются прочитанными

	_, err := ms.db.ExecContext(ctx, query,
		msg.ID,
		ourJID,
		msg.Chat,
		msg.From,
		msg.To,
		msg.Text,
		msg.Timestamp,
		"text", // message_type - можем расширить позже для других типов
		msg.QuotedMessageID,
		isFromMe, // is_from_me
		isRead,   // is_read - входящие сообщения непрочитанные, исходящие прочитанные
	)

	return err
}

// GetChatMessages retrieves messages for a specific chat with pagination
func (ms *MessageStore) GetChatMessages(ctx context.Context, ourJID, chatJID string, count int, beforeMessageID string) ([]types.Message, error) {
	var query string
	var args []interface{}

	if beforeMessageID != "" {
		// Get messages before a specific message (for pagination)
		query = `
			SELECT id, sender_jid, recipient_jid, message_text, timestamp, quoted_message_id
			FROM messages 
			WHERE our_jid = $1 AND chat_jid = $2 AND timestamp < (
				SELECT timestamp FROM messages WHERE id = $3 AND our_jid = $1
			)
			ORDER BY timestamp DESC 
			LIMIT $4
		`
		args = []interface{}{ourJID, chatJID, beforeMessageID, count}
	} else {
		// Get latest messages
		query = `
			SELECT id, sender_jid, recipient_jid, message_text, timestamp, quoted_message_id
			FROM messages 
			WHERE our_jid = $1 AND chat_jid = $2
			ORDER BY timestamp DESC 
			LIMIT $3
		`
		args = []interface{}{ourJID, chatJID, count}
	}

	rows, err := ms.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query messages: %w", err)
	}
	defer rows.Close()

	var messages []types.Message
	for rows.Next() {
		var msg types.Message
		var recipientJID sql.NullString
		var quotedMessageID sql.NullString

		err := rows.Scan(
			&msg.ID,
			&msg.From,
			&recipientJID,
			&msg.Text,
			&msg.Timestamp,
			&quotedMessageID,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan message: %w", err)
		}

		msg.Chat = chatJID
		if recipientJID.Valid {
			msg.To = recipientJID.String
		}
		if quotedMessageID.Valid {
			msg.QuotedMessageID = quotedMessageID.String
		}

		messages = append(messages, msg)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating messages: %w", err)
	}

	// Reverse the slice to get chronological order (oldest first)
	for i := len(messages)/2 - 1; i >= 0; i-- {
		opp := len(messages) - 1 - i
		messages[i], messages[opp] = messages[opp], messages[i]
	}

	return messages, nil
}

// GetAllMessages retrieves all messages for a user (used for counting)
func (ms *MessageStore) GetAllMessages(ctx context.Context, ourJID string) ([]types.Message, error) {
	query := `
		SELECT id, chat_jid, sender_jid, recipient_jid, message_text, timestamp, quoted_message_id
		FROM messages 
		WHERE our_jid = $1
		ORDER BY timestamp DESC
	`

	rows, err := ms.db.QueryContext(ctx, query, ourJID)
	if err != nil {
		return nil, fmt.Errorf("failed to query all messages: %w", err)
	}
	defer rows.Close()

	var messages []types.Message
	for rows.Next() {
		var msg types.Message
		var recipientJID sql.NullString
		var quotedMessageID sql.NullString

		err := rows.Scan(
			&msg.ID,
			&msg.Chat,
			&msg.From,
			&recipientJID,
			&msg.Text,
			&msg.Timestamp,
			&quotedMessageID,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan message: %w", err)
		}

		if recipientJID.Valid {
			msg.To = recipientJID.String
		}
		if quotedMessageID.Valid {
			msg.QuotedMessageID = quotedMessageID.String
		}

		messages = append(messages, msg)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating all messages: %w", err)
	}

	return messages, nil
}

// GetChatMessageCount returns the total count of messages in a chat
func (ms *MessageStore) GetChatMessageCount(ctx context.Context, ourJID, chatJID string) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM messages WHERE our_jid = $1 AND chat_jid = $2`

	err := ms.db.QueryRowContext(ctx, query, ourJID, chatJID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count messages: %w", err)
	}

	return count, nil
}

// DeleteOldMessages deletes messages older than the specified number of days
func (ms *MessageStore) DeleteOldMessages(ctx context.Context, ourJID string, daysToKeep int) (int64, error) {
	query := `
		DELETE FROM messages 
		WHERE our_jid = $1 AND created_at < NOW() - INTERVAL '%d days'
	`

	result, err := ms.db.ExecContext(ctx, fmt.Sprintf(query, daysToKeep), ourJID)
	if err != nil {
		return 0, fmt.Errorf("failed to delete old messages: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get rows affected: %w", err)
	}

	return rowsAffected, nil
}

// MarkMessagesAsRead marks messages as read in a specific chat
func (ms *MessageStore) MarkMessagesAsRead(ctx context.Context, ourJID, chatJID string) error {
	query := `UPDATE messages SET is_read = true WHERE our_jid = $1 AND chat_jid = $2 AND is_read = false`

	_, err := ms.db.ExecContext(ctx, query, ourJID, chatJID)
	if err != nil {
		return fmt.Errorf("failed to mark messages as read: %w", err)
	}

	return nil
}

// GetUnreadMessages retrieves unread messages with optional chat filter
func (ms *MessageStore) GetUnreadMessages(ctx context.Context, ourJID string, chatJID string, count int) ([]types.Message, error) {
	var query string
	var args []interface{}

	if chatJID != "" {
		// Get unread messages from a specific chat
		query = `
			SELECT id, sender_jid, recipient_jid, message_text, timestamp, quoted_message_id, chat_jid
			FROM messages 
			WHERE our_jid = $1 AND chat_jid = $2 AND is_read = false
			ORDER BY timestamp DESC 
			LIMIT $3
		`
		args = []interface{}{ourJID, chatJID, count}
	} else {
		// Get unread messages from all chats
		query = `
			SELECT id, sender_jid, recipient_jid, message_text, timestamp, quoted_message_id, chat_jid
			FROM messages 
			WHERE our_jid = $1 AND is_read = false
			ORDER BY timestamp DESC 
			LIMIT $2
		`
		args = []interface{}{ourJID, count}
	}

	rows, err := ms.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query unread messages: %w", err)
	}
	defer rows.Close()

	var messages []types.Message
	for rows.Next() {
		var msg types.Message
		var recipientJID sql.NullString
		var quotedMessageID sql.NullString

		err := rows.Scan(
			&msg.ID,
			&msg.From,
			&recipientJID,
			&msg.Text,
			&msg.Timestamp,
			&quotedMessageID,
			&msg.Chat,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan unread message: %w", err)
		}

		if recipientJID.Valid {
			msg.To = recipientJID.String
		}
		if quotedMessageID.Valid {
			msg.QuotedMessageID = quotedMessageID.String
		}

		messages = append(messages, msg)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating unread messages: %w", err)
	}

	// Reverse the slice to get chronological order (oldest first)
	for i := len(messages)/2 - 1; i >= 0; i-- {
		opp := len(messages) - 1 - i
		messages[i], messages[opp] = messages[opp], messages[i]
	}

	return messages, nil
}
