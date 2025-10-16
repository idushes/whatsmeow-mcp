package tools

import (
	"context"
	"time"
	"whatsmeow-mcp/internal/client"
	"whatsmeow-mcp/internal/qrcode"
	"whatsmeow-mcp/internal/types"

	"github.com/mark3labs/mcp-go/mcp"
)

// GetQRCodeTool creates and returns the get_qr_code MCP tool
func GetQRCodeTool(whatsappClient client.WhatsAppClientInterface) mcp.Tool {
	tool := mcp.NewTool("get_qr_code",
		mcp.WithDescription("Generate QR code for WhatsApp Web authentication."),
	)

	return tool
}

// HandleGetQRCode handles the get_qr_code tool execution
func HandleGetQRCode(whatsappClient client.WhatsAppClientInterface, qrGenerator *qrcode.QRCodeGenerator) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		qrCode := whatsappClient.GetQRCode()

		// Generate QR code image with base64
		qrResult, err := qrGenerator.GenerateQRCodeWithBase64(qrCode)
		if err != nil {
			errorResult := types.StandardResponse{
				Success: false,
				Error: &types.ErrorInfo{
					Code:    "QR_GENERATION_FAILED",
					Message: "Failed to generate QR code image",
					Details: err.Error(),
				},
			}
			return mcp.NewToolResultStructured(errorResult, "Failed to generate QR code"), nil
		}

		expiresAt := time.Now().Add(30 * time.Second).Unix()

		result := types.QRCodeResponse{
			QRCode:    qrCode,
			Code:      qrCode,
			ImageURL:  qrResult.ImageURL,
			Timeout:   30,
			Success:   true,
			ExpiresAt: expiresAt,
		}

		// Create content with text, image, and resource link
		content := []mcp.Content{
			mcp.NewTextContent("QR code generated successfully. Expires in 30 seconds. Scan with WhatsApp to login."),
			mcp.NewImageContent(qrResult.Base64, "image/png"),
			mcp.NewResourceLink(
				qrResult.ImageURL,
				"WhatsApp QR Code",
				"Scan this QR code with WhatsApp to login",
				"image/png",
			),
		}

		return &mcp.CallToolResult{
			Content:           content,
			StructuredContent: result,
			IsError:           false,
		}, nil
	}
}
