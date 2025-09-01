package tools

import (
	"context"
	"encoding/json"
	"time"
	"whatsmeow-mcp/internal/client"
	"whatsmeow-mcp/internal/qrcode"
	"whatsmeow-mcp/internal/types"

	"github.com/mark3labs/mcp-go/mcp"
)

// GetQRCodeTool creates and returns the get_qr_code MCP tool
func GetQRCodeTool(whatsappClient client.WhatsAppClientInterface) mcp.Tool {
	tool := mcp.NewTool("get_qr_code",
		mcp.WithDescription("Generate QR code for WhatsApp Web login. Returns both the raw QR code string and a hosted image URL that can be displayed directly to users for scanning with WhatsApp mobile app. The QR code expires after 30 seconds and the image file is automatically cleaned up after 5 minutes."),
	)

	return tool
}

// HandleGetQRCode handles the get_qr_code tool execution
func HandleGetQRCode(whatsappClient client.WhatsAppClientInterface, qrGenerator *qrcode.QRCodeGenerator) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		qrCode := whatsappClient.GetQRCode()

		// Generate QR code image
		imageURL, err := qrGenerator.GenerateQRCode(qrCode)
		if err != nil {
			errorResult := types.StandardResponse{
				Success: false,
				Error: &types.ErrorInfo{
					Code:    "QR_GENERATION_FAILED",
					Message: "Failed to generate QR code image",
					Details: err.Error(),
				},
			}
			content, _ := json.Marshal(errorResult)
			return &mcp.CallToolResult{
				Content: []mcp.Content{
					mcp.NewTextContent(string(content)),
				},
			}, nil
		}

		expiresAt := time.Now().Add(30 * time.Second).Unix()

		result := types.QRCodeResponse{
			QRCode:    qrCode,
			Code:      qrCode,
			ImageURL:  imageURL,
			Timeout:   30,
			Success:   true,
			ExpiresAt: expiresAt,
		}

		// QR code expires after 30 seconds
		go func() {
			time.Sleep(30 * time.Second)
			whatsappClient.UpdateQRCode()
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
