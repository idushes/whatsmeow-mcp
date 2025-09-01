package tools

import (
	"context"
	"encoding/json"
	"whatsmeow-mcp/internal/client"
	"whatsmeow-mcp/internal/types"

	"github.com/mark3labs/mcp-go/mcp"
)

// GetUnreadMessagesTool creates and returns the get_unread_messages MCP tool
func GetUnreadMessagesTool(whatsappClient client.WhatsAppClientInterface) mcp.Tool {
	tool := mcp.NewTool("get_unread_messages",
		mcp.WithDescription("Get unread messages from WhatsApp. Can retrieve unread messages from all chats or filter by specific chat. Returns messages that have not been marked as read in the database."),
		mcp.WithInputSchema[types.GetUnreadMessagesParams](),
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
			content, _ := json.Marshal(result)
			return &mcp.CallToolResult{
				Content: []mcp.Content{
					mcp.NewTextContent(string(content)),
				},
			}, nil
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
		unreadMessages := whatsappClient.GetUnreadMessages(params.Chat, params.Count)

		result := types.UnreadMessagesResponse{
			Messages: unreadMessages,
			Success:  true,
			Chat:     params.Chat,
			Count:    len(unreadMessages),
		}

		content, err := json.Marshal(result)
		if err != nil {
			errorResult := types.StandardResponse{
				Success: false,
				Error: &types.ErrorInfo{
					Code:    "MARSHAL_ERROR",
					Message: "Failed to serialize response",
					Details: err.Error(),
				},
			}
			content, _ = json.Marshal(errorResult)
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.NewTextContent(string(content)),
			},
		}, nil
	}
}
