package tools

import (
	"context"
	"encoding/json"
	"time"
	"whatsmeow-mcp/internal/client"
	"whatsmeow-mcp/internal/types"

	"github.com/mark3labs/mcp-go/mcp"
)

// GetQRCodeTool creates and returns the get_qr_code MCP tool
func GetQRCodeTool(client *client.WhatsAppClient) mcp.Tool {
	tool := mcp.NewTool("get_qr_code",
		mcp.WithDescription("Generate QR code for WhatsApp Web login. Returns a QR code string that can be scanned with WhatsApp mobile app to authenticate the session. The QR code expires after 30 seconds and needs to be refreshed."),
		mcp.WithInputSchema[types.DummyParams](),
	)

	return tool
}

// HandleGetQRCode handles the get_qr_code tool execution
func HandleGetQRCode(client *client.WhatsAppClient) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		qrCode := client.GetQRCode()

		result := types.QRCodeResponse{
			QRCode:  qrCode,
			Code:    qrCode,
			Timeout: 30,
			Success: true,
		}

		// QR code expires after 30 seconds
		go func() {
			time.Sleep(30 * time.Second)
			client.UpdateQRCode()
		}()

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
