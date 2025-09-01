package tools

import (
	"context"
	"encoding/json"
	"whatsmeow-mcp/internal/client"
	"whatsmeow-mcp/internal/types"

	"github.com/mark3labs/mcp-go/mcp"
)

// GetChatHistoryTool creates and returns the get_chat_history MCP tool
func GetChatHistoryTool(whatsappClient client.WhatsAppClientInterface) mcp.Tool {
	tool := mcp.NewTool("get_chat_history",
		mcp.WithDescription("Retrieve message history from a WhatsApp conversation with pagination support."),
		mcp.WithString("chat",
			mcp.Required(),
			mcp.Description("WhatsApp JID (chat identifier) to retrieve messages from"),
		),
		mcp.WithNumber("count",
			mcp.Description("Maximum number of messages to retrieve (default: 50, max: 100)"),
		),
		mcp.WithString("before_message_id",
			mcp.Description("Optional message ID to retrieve messages before this point (for pagination)"),
		),
	)

	return tool
}

// HandleGetChatHistory handles the get_chat_history tool execution
func HandleGetChatHistory(whatsappClient client.WhatsAppClientInterface) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		var params types.GetChatHistoryParams
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
					Code:    "MISSING_PARAMETERS",
					Message: "Required parameter 'chat' (WhatsApp JID) must be provided",
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

		// Retrieve messages for the specific chat
		chatMessages := whatsappClient.GetChatMessages(params.Chat, params.Count, params.BeforeMessageID)

		// For now, we'll determine hasMore by checking if we got the full requested count
		// In a more sophisticated implementation, we could add a method to get total count
		hasMore := len(chatMessages) == params.Count

		result := types.ChatHistoryResponse{
			Messages: chatMessages,
			HasMore:  hasMore,
			Success:  true,
			Chat:     params.Chat,
			Count:    len(chatMessages),
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
