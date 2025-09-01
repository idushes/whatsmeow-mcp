package tools

import (
	"log"
	"whatsmeow-mcp/internal/client"
	"whatsmeow-mcp/internal/qrcode"

	"github.com/mark3labs/mcp-go/server"
)

// RegisterAllTools registers all available WhatsApp MCP tools with the server
func RegisterAllTools(mcpServer *server.MCPServer, whatsappClient client.WhatsAppClientInterface, qrGenerator *qrcode.QRCodeGenerator) {
	// Register is_logged_in tool
	isLoggedInTool := IsLoggedInTool(whatsappClient)
	mcpServer.AddTool(isLoggedInTool, HandleIsLoggedIn(whatsappClient))

	// Register get_qr_code tool
	getQRCodeTool := GetQRCodeTool(whatsappClient)
	mcpServer.AddTool(getQRCodeTool, HandleGetQRCode(whatsappClient, qrGenerator))

	// Register send_message tool
	sendMessageTool := SendMessageTool(whatsappClient)
	mcpServer.AddTool(sendMessageTool, HandleSendMessage(whatsappClient))

	// Register is_on_whatsapp tool
	isOnWhatsappTool := IsOnWhatsappTool(whatsappClient)
	mcpServer.AddTool(isOnWhatsappTool, HandleIsOnWhatsapp(whatsappClient))

	// Register get_chat_history tool
	getChatHistoryTool := GetChatHistoryTool(whatsappClient)
	mcpServer.AddTool(getChatHistoryTool, HandleGetChatHistory(whatsappClient))

	log.Println("Successfully registered 5 WhatsApp MCP tools:")
	log.Println("  - is_logged_in: Check authentication status")
	log.Println("  - get_qr_code: Generate QR code for login")
	log.Println("  - send_message: Send text messages")
	log.Println("  - is_on_whatsapp: Check phone number registration")
	log.Println("  - get_chat_history: Retrieve chat message history")
}
