package tools

import (
	"context"
	"encoding/json"
	"whatsmeow-mcp/internal/client"
	"whatsmeow-mcp/internal/types"

	"github.com/mark3labs/mcp-go/mcp"
)

// SendMessageTool creates and returns the send_message MCP tool
func SendMessageTool(whatsappClient client.WhatsAppClientInterface) mcp.Tool {
	tool := mcp.NewTool("send_message",
		mcp.WithDescription("Send a text message to a WhatsApp chat or contact. Requires the user to be authenticated (logged in) first.\n\nParameters:\n- to: WhatsApp JID (recipient identifier) in format 'phonenumber@s.whatsapp.net' (e.g., '1234567890@s.whatsapp.net') or group JID ending with '@g.us'\n- text: The message content to send (plain text)\n- quoted_message_id: (Optional) ID of a previous message to reply to/quote\n\nExample usage:\n- Send to phone: to='1234567890@s.whatsapp.net', text='Hello!'\n- Reply to message: to='1234567890@s.whatsapp.net', text='Thanks!', quoted_message_id='msg_123'"),
		mcp.WithInputSchema[types.SendMessageParams](),
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
			content, _ := json.Marshal(result)
			return &mcp.CallToolResult{
				Content: []mcp.Content{
					mcp.NewTextContent(string(content)),
				},
			}, nil
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

		// Send message using client interface
		response, err := whatsappClient.SendMessage(params.To, params.Text, params.QuotedMessageID)
		if err != nil {
			result := types.StandardResponse{
				Success: false,
				Error: &types.ErrorInfo{
					Code:    "SEND_FAILED",
					Message: "Failed to send message",
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

		result := response

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
