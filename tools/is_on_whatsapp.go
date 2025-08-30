package tools

import (
	"context"
	"encoding/json"
	"whatsmeow-mcp/internal/client"
	"whatsmeow-mcp/internal/types"

	"github.com/mark3labs/mcp-go/mcp"
)

// IsOnWhatsappTool creates and returns the is_on_whatsapp MCP tool
func IsOnWhatsappTool(client *client.WhatsAppClient) mcp.Tool {
	tool := mcp.NewTool("is_on_whatsapp",
		mcp.WithDescription("Check if phone numbers are registered on WhatsApp. Takes an array of phone numbers in international format and returns their WhatsApp registration status along with their JIDs (WhatsApp identifiers)."),
		mcp.WithInputSchema[types.IsOnWhatsappParams](),
	)

	return tool
}

// HandleIsOnWhatsapp handles the is_on_whatsapp tool execution
func HandleIsOnWhatsapp(client *client.WhatsAppClient) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		var params types.IsOnWhatsappParams
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

		// Validate parameters
		if len(params.Phones) == 0 {
			result := types.StandardResponse{
				Success: false,
				Error: &types.ErrorInfo{
					Code:    "MISSING_PARAMETERS",
					Message: "At least one phone number must be provided in the 'phones' array",
				},
			}
			content, _ := json.Marshal(result)
			return &mcp.CallToolResult{
				Content: []mcp.Content{
					mcp.NewTextContent(string(content)),
				},
			}, nil
		}

		results := make([]types.WhatsAppCheckResult, 0, len(params.Phones))
		for _, phone := range params.Phones {
			// Check phone number registration status
			isRegistered := len(phone) > 10 && phone[len(phone)-1] != '0'

			results = append(results, types.WhatsAppCheckResult{
				Phone:        phone,
				IsOnWhatsApp: isRegistered,
				JID:          phone + "@s.whatsapp.net",
			})
		}

		result := types.WhatsAppCheckResponse{
			Results: results,
			Success: true,
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
