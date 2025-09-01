package tools

import (
	"context"
	"encoding/json"
	"whatsmeow-mcp/internal/client"
	"whatsmeow-mcp/internal/types"

	"github.com/mark3labs/mcp-go/mcp"
)

// MarkMessagesAsReadTool creates and returns the mark_messages_as_read MCP tool
func MarkMessagesAsReadTool(whatsappClient client.WhatsAppClientInterface) mcp.Tool {
	tool := mcp.NewTool("mark_messages_as_read",
		mcp.WithDescription("Mark all unread messages in a specific chat as read."),
		mcp.WithString("chat",
			mcp.Description("WhatsApp JID (chat identifier) to mark messages as read in this chat. For phone numbers: 'phonenumber@s.whatsapp.net' (e.g. '1234567890@s.whatsapp.net'). For groups: 'groupid@g.us'"),
		),
	)

	return tool
}

// HandleMarkMessagesAsRead handles the mark_messages_as_read tool execution
func HandleMarkMessagesAsRead(whatsappClient client.WhatsAppClientInterface) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		var params types.MarkMessagesAsReadParams
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

		// Validate required parameters
		if params.Chat == "" {
			result := types.StandardResponse{
				Success: false,
				Error: &types.ErrorInfo{
					Code:    "MISSING_CHAT_PARAMETER",
					Message: "Chat parameter is required",
					Details: "Please provide a valid WhatsApp JID (chat identifier)",
				},
			}
			content, _ := json.Marshal(result)
			return &mcp.CallToolResult{
				Content: []mcp.Content{
					mcp.NewTextContent(string(content)),
				},
			}, nil
		}

		// Mark messages as read
		err := whatsappClient.MarkMessagesAsRead(params.Chat)
		if err != nil {
			result := types.StandardResponse{
				Success: false,
				Error: &types.ErrorInfo{
					Code:    "MARK_AS_READ_FAILED",
					Message: "Failed to mark messages as read",
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

		// Success response
		result := types.MarkMessagesAsReadResponse{
			Success: true,
			Chat:    params.Chat,
			Message: "All unread messages in the chat have been marked as read",
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
