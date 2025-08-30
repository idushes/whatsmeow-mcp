package tools

import (
	"context"
	"encoding/json"
	"whatsmeow-mcp/internal/client"
	"whatsmeow-mcp/internal/types"

	"github.com/mark3labs/mcp-go/mcp"
)

// IsLoggedInTool creates and returns the is_logged_in MCP tool
func IsLoggedInTool(client *client.WhatsAppClient) mcp.Tool {
	tool := mcp.NewTool("is_logged_in",
		mcp.WithDescription("Check if user is authenticated with WhatsApp. Returns the current login status of the WhatsApp client session."),
		mcp.WithInputSchema[types.DummyParams](),
	)

	return tool
}

// HandleIsLoggedIn handles the is_logged_in tool execution
func HandleIsLoggedIn(client *client.WhatsAppClient) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		result := types.LoginStatusResponse{
			LoggedIn: client.IsLoggedIn(),
			Success:  true,
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
