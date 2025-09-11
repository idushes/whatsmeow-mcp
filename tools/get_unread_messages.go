package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"whatsmeow-mcp/internal/client"
	"whatsmeow-mcp/internal/types"

	"github.com/mark3labs/mcp-go/mcp"
)

// GetUnreadMessagesTool creates and returns the get_unread_messages MCP tool
func GetUnreadMessagesTool(whatsappClient client.WhatsAppClientInterface) mcp.Tool {
	tool := mcp.NewTool("get_unread_messages",
		mcp.WithDescription("Retrieve unread messages from WhatsApp chats."),
		mcp.WithString("chat",
			mcp.Description("Optional WhatsApp JID to filter unread messages from a specific chat. If omitted, returns unread messages from all chats"),
		),
		mcp.WithNumber("count",
			mcp.Description("Maximum number of unread messages to retrieve (default: 50, max: 100)"),
		),
	)

	return tool
}

// HandleGetUnreadMessages handles the get_unread_messages tool execution
func HandleGetUnreadMessages(whatsappClient client.WhatsAppClientInterface) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		var params types.GetUnreadMessagesParams
		argumentsBytes, _ := json.Marshal(request.Params.Arguments)
		if err := json.Unmarshal(argumentsBytes, &params); err != nil {
			result := types.StandardResponse{
				Success: false,
				Error: &types.ErrorInfo{
					Code:    "INVALID_PARAMETERS",
					Message: "Failed to parse parameters",
					Details: err.Error(),
				},
			}
			return mcp.NewToolResultStructured(result, "Failed to parse parameters"), nil
		}

		// Set default count if not provided or invalid
		if params.Count <= 0 {
			params.Count = 50
		}

		// Limit maximum count to prevent excessive data retrieval
		if params.Count > 100 {
			params.Count = 100
		}

		// Retrieve unread messages
		allUnreadMessages := whatsappClient.GetUnreadMessages(params.Chat, params.Count)

		// Filter out empty messages
		var unreadMessages []types.Message
		for _, msg := range allUnreadMessages {
			// Skip messages with empty text content and no meaningful data
			// Only include messages that have actual text content
			if msg.Text != "" && len(strings.TrimSpace(msg.Text)) > 0 {
				unreadMessages = append(unreadMessages, msg)
			}
		}

		result := types.UnreadMessagesResponse{
			Messages: unreadMessages,
			Success:  true,
			Chat:     params.Chat,
			Count:    len(unreadMessages),
		}

		// Create fallback text for backward compatibility
		fallbackText := fmt.Sprintf("Found %d unread messages", len(unreadMessages))
		if params.Chat != "" {
			fallbackText += fmt.Sprintf(" in chat %s", params.Chat)
		}

		return mcp.NewToolResultStructured(result, fallbackText), nil
	}
}
