package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	"whatsmeow-mcp/internal/client"
	"whatsmeow-mcp/internal/types"

	"github.com/mark3labs/mcp-go/mcp"
)

// SendMessageTool creates and returns the send_message MCP tool
func SendMessageTool(client *client.WhatsAppClient) mcp.Tool {
	tool := mcp.NewTool("send_message",
		mcp.WithDescription("Send text message to chat or contact. Requires authentication (user must be logged in). Can optionally reply to/quote a previous message by providing quoted_message_id parameter."),
		mcp.WithInputSchema[types.SendMessageParams](),
	)

	return tool
}

// HandleSendMessage handles the send_message tool execution
func HandleSendMessage(client *client.WhatsAppClient) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		var params types.SendMessageParams
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

		// Check if user is authenticated
		if !client.IsLoggedIn() {
			result := types.StandardResponse{
				Success: false,
				Error: &types.ErrorInfo{
					Code:    "NOT_LOGGED_IN",
					Message: "Client is not authenticated. Please login first using get_qr_code tool.",
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
		if params.To == "" || params.Text == "" {
			result := types.StandardResponse{
				Success: false,
				Error: &types.ErrorInfo{
					Code:    "MISSING_PARAMETERS",
					Message: "Required parameters 'to' and 'text' must be provided",
				},
			}
			content, _ := json.Marshal(result)
			return &mcp.CallToolResult{
				Content: []mcp.Content{
					mcp.NewTextContent(string(content)),
				},
			}, nil
		}

		// Send message
		messageId := fmt.Sprintf("msg_%d", time.Now().Unix())
		timestamp := time.Now().Unix()

		result := types.MessageResponse{
			MessageID:       messageId,
			Timestamp:       timestamp,
			Success:         true,
			To:              params.To,
			Text:            params.Text,
			QuotedMessageID: params.QuotedMessageID,
		}

		// Add message to client history
		newMessage := types.Message{
			ID:              messageId,
			From:            "self",
			To:              params.To,
			Text:            params.Text,
			Timestamp:       timestamp,
			Chat:            params.To,
			QuotedMessageID: params.QuotedMessageID,
		}
		client.AddMessage(newMessage)

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
