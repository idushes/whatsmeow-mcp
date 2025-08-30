package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// Config holds the server configuration
type Config struct {
	Port          int
	Host          string
	ServerName    string
	ServerVersion string
	LogLevel      string
}

// loadConfig loads configuration from environment variables
func loadConfig() *Config {
	// Try to load .env file, ignore error if it doesn't exist
	_ = godotenv.Load("config.env")

	config := &Config{
		Port:          3000,
		Host:          "localhost",
		ServerName:    "whatsmeow-mcp",
		ServerVersion: "1.0.0",
		LogLevel:      "info",
	}

	if port := os.Getenv("PORT"); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			config.Port = p
		}
	}

	if host := os.Getenv("HOST"); host != "" {
		config.Host = host
	}

	if name := os.Getenv("SERVER_NAME"); name != "" {
		config.ServerName = name
	}

	if version := os.Getenv("SERVER_VERSION"); version != "" {
		config.ServerVersion = version
	}

	if level := os.Getenv("LOG_LEVEL"); level != "" {
		config.LogLevel = level
	}

	return config
}

// WhatsAppClient simulates a WhatsApp client state
type WhatsAppClient struct {
	isLoggedIn bool
	qrCode     string
	messages   []map[string]interface{}
}

// Tool parameter structures
type SendMessageParams struct {
	To              string `json:"to" description:"Recipient JID"`
	Text            string `json:"text" description:"Message text content"`
	QuotedMessageID string `json:"quoted_message_id,omitempty" description:"ID of message to quote/reply to"`
}

type IsOnWhatsappParams struct {
	Phones []string `json:"phones" description:"Phone numbers in international format"`
}

type GetChatHistoryParams struct {
	Chat            string `json:"chat" description:"Chat JID"`
	Count           int    `json:"count,omitempty" description:"Number of messages to retrieve (default: 50)"`
	BeforeMessageID string `json:"before_message_id,omitempty" description:"Get messages before this ID"`
}

// Global client instance for simulation
var client = &WhatsAppClient{
	isLoggedIn: false,
	qrCode:     "2@ABC123DEF456GHI789JKL012MNO345PQR678STU901VWX234YZ,567890abcdef1234567890abcdef1234567890ab,cdefghijklmnopqrstuvwxyz",
	messages: []map[string]interface{}{
		{
			"id":        "msg_001",
			"from":      "1234567890@s.whatsapp.net",
			"text":      "Hello! This is a test message.",
			"timestamp": time.Now().Unix() - 3600,
			"chat":      "1234567890@s.whatsapp.net",
		},
		{
			"id":        "msg_002",
			"from":      "9876543210@s.whatsapp.net",
			"text":      "How are you doing?",
			"timestamp": time.Now().Unix() - 1800,
			"chat":      "9876543210@s.whatsapp.net",
		},
	},
}

func main() {
	config := loadConfig()

	log.Printf("Starting %s v%s", config.ServerName, config.ServerVersion)
	log.Printf("Configuration: Host=%s, Port=%d", config.Host, config.Port)

	// Create MCP server
	mcpServer := server.NewMCPServer(
		config.ServerName,
		config.ServerVersion,
	)

	// Register tools
	registerTools(mcpServer)

	// Check if running in stdio mode
	if len(os.Args) > 1 && os.Args[1] == "stdio" {
		log.Println("Starting MCP server in stdio mode")
		err := server.ServeStdio(mcpServer)
		if err != nil {
			log.Fatal("Failed to start stdio server:", err)
		}
	} else {
		log.Printf("Starting MCP server in SSE mode on %s:%d", config.Host, config.Port)
		sseServer := server.NewSSEServer(mcpServer,
			server.WithSSEEndpoint("/sse"),
		)
		err := sseServer.Start(fmt.Sprintf("%s:%d", config.Host, config.Port))
		if err != nil {
			log.Fatal("Failed to start SSE server:", err)
		}
	}
}

func registerTools(mcpServer *server.MCPServer) {
	// Tool: is_logged_in
	isLoggedInTool := mcp.NewTool("is_logged_in",
		mcp.WithDescription("Check if user is authenticated"),
	)
	mcpServer.AddTool(isLoggedInTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		result := map[string]interface{}{
			"logged_in": client.isLoggedIn,
			"success":   true,
		}

		content, _ := json.Marshal(result)
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.NewTextContent(string(content)),
			},
		}, nil
	})

	// Tool: get_qr_code
	getQrCodeTool := mcp.NewTool("get_qr_code",
		mcp.WithDescription("Generate QR code for WhatsApp Web login"),
	)
	mcpServer.AddTool(getQrCodeTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		result := map[string]interface{}{
			"qr_code": client.qrCode,
			"code":    client.qrCode,
			"timeout": 30,
			"success": true,
		}

		// Simulate QR code expiration after some time
		go func() {
			time.Sleep(30 * time.Second)
			client.qrCode = "2@XYZ789ABC012DEF345GHI678JKL901MNO234PQR567STU890VWX123YZ,456789abcdef0123456789abcdef0123456789ab,cdefghijklmnopqrstuvwxyz"
		}()

		content, _ := json.Marshal(result)
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.NewTextContent(string(content)),
			},
		}, nil
	})

	// Tool: send_message
	sendMessageTool := mcp.NewTool("send_message",
		mcp.WithDescription("Send text message to chat or contact"),
		mcp.WithInputSchema[SendMessageParams](),
	)
	mcpServer.AddTool(sendMessageTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		var params SendMessageParams
		argumentsBytes, _ := json.Marshal(request.Params.Arguments)
		if err := json.Unmarshal(argumentsBytes, &params); err != nil {
			result := map[string]interface{}{
				"success": false,
				"error": map[string]interface{}{
					"code":    "INVALID_PARAMETERS",
					"message": "Failed to parse parameters",
				},
			}
			content, _ := json.Marshal(result)
			return &mcp.CallToolResult{
				Content: []mcp.Content{
					mcp.NewTextContent(string(content)),
				},
			}, nil
		}

		if !client.isLoggedIn {
			result := map[string]interface{}{
				"success": false,
				"error": map[string]interface{}{
					"code":    "NOT_LOGGED_IN",
					"message": "Client is not authenticated",
				},
			}
			content, _ := json.Marshal(result)
			return &mcp.CallToolResult{
				Content: []mcp.Content{
					mcp.NewTextContent(string(content)),
				},
			}, nil
		}

		// Simulate message sending
		messageId := fmt.Sprintf("msg_%d", time.Now().Unix())
		timestamp := time.Now().Unix()

		result := map[string]interface{}{
			"message_id":        messageId,
			"timestamp":         timestamp,
			"success":           true,
			"to":                params.To,
			"text":              params.Text,
			"quoted_message_id": params.QuotedMessageID,
		}

		// Add message to history
		newMessage := map[string]interface{}{
			"id":        messageId,
			"from":      "self",
			"to":        params.To,
			"text":      params.Text,
			"timestamp": timestamp,
			"chat":      params.To,
		}
		if params.QuotedMessageID != "" {
			newMessage["quoted_message_id"] = params.QuotedMessageID
		}
		client.messages = append(client.messages, newMessage)

		content, _ := json.Marshal(result)
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.NewTextContent(string(content)),
			},
		}, nil
	})

	// Tool: is_on_whatsapp
	isOnWhatsappTool := mcp.NewTool("is_on_whatsapp",
		mcp.WithDescription("Check if phone numbers are registered on WhatsApp"),
		mcp.WithInputSchema[IsOnWhatsappParams](),
	)
	mcpServer.AddTool(isOnWhatsappTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		var params IsOnWhatsappParams
		argumentsBytes, _ := json.Marshal(request.Params.Arguments)
		if err := json.Unmarshal(argumentsBytes, &params); err != nil {
			result := map[string]interface{}{
				"success": false,
				"error": map[string]interface{}{
					"code":    "INVALID_PARAMETERS",
					"message": "Failed to parse parameters",
				},
			}
			content, _ := json.Marshal(result)
			return &mcp.CallToolResult{
				Content: []mcp.Content{
					mcp.NewTextContent(string(content)),
				},
			}, nil
		}

		results := make([]map[string]interface{}, 0, len(params.Phones))
		for _, phone := range params.Phones {
			// Simulate random registration status
			isRegistered := len(phone) > 10 && phone[len(phone)-1] != '0'

			results = append(results, map[string]interface{}{
				"phone":          phone,
				"is_on_whatsapp": isRegistered,
				"jid":            phone + "@s.whatsapp.net",
			})
		}

		result := map[string]interface{}{
			"results": results,
			"success": true,
		}

		content, _ := json.Marshal(result)
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.NewTextContent(string(content)),
			},
		}, nil
	})

	// Tool: get_chat_history
	getChatHistoryTool := mcp.NewTool("get_chat_history",
		mcp.WithDescription("Get chat message history"),
		mcp.WithInputSchema[GetChatHistoryParams](),
	)
	mcpServer.AddTool(getChatHistoryTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		var params GetChatHistoryParams
		argumentsBytes, _ := json.Marshal(request.Params.Arguments)
		if err := json.Unmarshal(argumentsBytes, &params); err != nil {
			result := map[string]interface{}{
				"success": false,
				"error": map[string]interface{}{
					"code":    "INVALID_PARAMETERS",
					"message": "Failed to parse parameters",
				},
			}
			content, _ := json.Marshal(result)
			return &mcp.CallToolResult{
				Content: []mcp.Content{
					mcp.NewTextContent(string(content)),
				},
			}, nil
		}

		if params.Count <= 0 {
			params.Count = 50
		}

		// Filter messages for the specific chat
		var chatMessages []map[string]interface{}
		for _, msg := range client.messages {
			if msgChat, ok := msg["chat"].(string); ok && msgChat == params.Chat {
				if params.BeforeMessageID == "" || msg["id"].(string) != params.BeforeMessageID {
					chatMessages = append(chatMessages, msg)
				}
			}
		}

		// Limit results
		if len(chatMessages) > params.Count {
			chatMessages = chatMessages[:params.Count]
		}

		result := map[string]interface{}{
			"messages": chatMessages,
			"has_more": len(client.messages) > params.Count,
			"success":  true,
			"chat":     params.Chat,
			"count":    len(chatMessages),
		}

		content, _ := json.Marshal(result)
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.NewTextContent(string(content)),
			},
		}, nil
	})

	log.Println("Registered 5 fake WhatsApp tools:")
	log.Println("  - is_logged_in")
	log.Println("  - get_qr_code")
	log.Println("  - send_message")
	log.Println("  - is_on_whatsapp")
	log.Println("  - get_chat_history")
}
