package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"whatsmeow-mcp/internal/client"
	"whatsmeow-mcp/internal/types"

	"github.com/mark3labs/mcp-go/mcp"
)

// IsOnWhatsappTool creates and returns the is_on_whatsapp MCP tool
func IsOnWhatsappTool(whatsappClient client.WhatsAppClientInterface) mcp.Tool {
	tool := mcp.NewTool("is_on_whatsapp",
		mcp.WithDescription("Check if phone numbers are registered on WhatsApp and get their JIDs."),
		mcp.WithArray("phones",
			mcp.Required(),
			mcp.Description("Array of phone numbers in international format (e.g., +1234567890) to check"),
			mcp.WithStringItems(mcp.Description("Phone number in international format")),
		),
	)

	return tool
}

// HandleIsOnWhatsapp handles the is_on_whatsapp tool execution
func HandleIsOnWhatsapp(whatsappClient client.WhatsAppClientInterface) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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
			return mcp.NewToolResultStructured(result, "Failed to parse parameters"), nil
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
			return mcp.NewToolResultStructured(result, "No phone numbers provided"), nil
		}

		// Use client interface to check phone numbers
		results, err := whatsappClient.IsOnWhatsApp(params.Phones)
		if err != nil {
			result := types.StandardResponse{
				Success: false,
				Error: &types.ErrorInfo{
					Code:    "CHECK_FAILED",
					Message: "Failed to check WhatsApp registration status",
					Details: err.Error(),
				},
			}
			return mcp.NewToolResultStructured(result, "Failed to check WhatsApp registration status"), nil
		}

		result := types.WhatsAppCheckResponse{
			Results: results,
			Success: true,
		}

		// Create fallback text for backward compatibility
		fallbackText := fmt.Sprintf("Checked %d phone numbers for WhatsApp registration", len(params.Phones))

		return mcp.NewToolResultStructured(result, fallbackText), nil
	}
}
