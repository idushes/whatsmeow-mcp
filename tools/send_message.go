package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"whatsmeow-mcp/internal/client"
	"whatsmeow-mcp/internal/types"

	"github.com/mark3labs/mcp-go/mcp"
)

// SendMessageTool creates and returns the send_message MCP tool
func SendMessageTool(whatsappClient client.WhatsAppClientInterface) mcp.Tool {
	tool := mcp.NewTool("send_message",
		mcp.WithDescription("Send a text message to a WhatsApp chat or contact. Requires authentication. IMPORTANT: When you send a message to a contact, your session will be automatically subscribed to receive real-time MCP notifications for all incoming messages from that contact. This means you'll receive 'notifications/message' events whenever the contact replies or sends new messages. Subscriptions are maintained per MCP session and prevent duplicate notifications."),
		mcp.WithString("to",
			mcp.Required(),
			mcp.Description("WhatsApp JID (recipient identifier) in format 'phonenumber@s.whatsapp.net' (e.g., '1234567890@s.whatsapp.net') or group JID ending with '@g.us'"),
		),
		mcp.WithString("text",
			mcp.Required(),
			mcp.Description("The message content to send (plain text)"),
		),
		mcp.WithString("quoted_message_id",
			mcp.Description("Optional ID of a previous message to reply to/quote"),
		),
	)

	return tool
}

// HandleSendMessage handles the send_message tool execution
func HandleSendMessage(whatsappClient client.WhatsAppClientInterface) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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
			return mcp.NewToolResultStructured(result, "Failed to parse parameters"), nil
		}

		// Check if user is authenticated
		if !whatsappClient.IsLoggedIn() {
			result := types.StandardResponse{
				Success: false,
				Error: &types.ErrorInfo{
					Code:    "NOT_LOGGED_IN",
					Message: "Client is not authenticated. Please login first using get_qr_code tool.",
				},
			}
			return mcp.NewToolResultStructured(result, "Not authenticated. Please login first."), nil
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
			return mcp.NewToolResultStructured(result, "Missing required parameters: 'to' and 'text'"), nil
		}

		// Send message using client interface (context contains session for auto-subscription)
		response, err := whatsappClient.SendMessage(ctx, params.To, params.Text, params.QuotedMessageID)
		if err != nil {
			result := types.StandardResponse{
				Success: false,
				Error: &types.ErrorInfo{
					Code:    "SEND_FAILED",
					Message: "Failed to send message",
					Details: err.Error(),
				},
			}
			return mcp.NewToolResultStructured(result, "Failed to send message"), nil
		}

		result := response

		// Create fallback text for backward compatibility
		fallbackText := fmt.Sprintf("Message sent successfully to %s. You are now subscribed to notifications from this chat.", params.To)

		return mcp.NewToolResultStructured(result, fallbackText), nil
	}
}
