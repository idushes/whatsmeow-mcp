package tools

import (
	"context"
	"whatsmeow-mcp/internal/client"
	"whatsmeow-mcp/internal/types"

	"github.com/mark3labs/mcp-go/mcp"
)

// IsLoggedInTool creates and returns the is_logged_in MCP tool
func IsLoggedInTool(whatsappClient client.WhatsAppClientInterface) mcp.Tool {
	tool := mcp.NewTool("is_logged_in",
		mcp.WithDescription("Check WhatsApp authentication status."),
	)

	return tool
}

// HandleIsLoggedIn handles the is_logged_in tool execution
func HandleIsLoggedIn(whatsappClient client.WhatsAppClientInterface) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		result := types.LoginStatusResponse{
			LoggedIn: whatsappClient.IsLoggedIn(),
			Success:  true,
		}

		// Create fallback text for backward compatibility
		var fallbackText string
		if result.LoggedIn {
			fallbackText = "Client is logged in to WhatsApp"
		} else {
			fallbackText = "Client is not logged in to WhatsApp"
		}

		return mcp.NewToolResultStructured(result, fallbackText), nil
	}
}
