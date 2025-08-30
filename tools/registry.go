package tools

import (
	"log"
	"whatsmeow-mcp/internal/client"

	"github.com/mark3labs/mcp-go/server"
)

// RegisterAllTools registers all available WhatsApp MCP tools with the server
func RegisterAllTools(mcpServer *server.MCPServer, client *client.WhatsAppClient) {
	// Register is_logged_in tool
	isLoggedInTool := IsLoggedInTool(client)
	mcpServer.AddTool(isLoggedInTool, HandleIsLoggedIn(client))

	// Register get_qr_code tool
	getQRCodeTool := GetQRCodeTool(client)
	mcpServer.AddTool(getQRCodeTool, HandleGetQRCode(client))

	// Register send_message tool
	sendMessageTool := SendMessageTool(client)
	mcpServer.AddTool(sendMessageTool, HandleSendMessage(client))

	// Register is_on_whatsapp tool
	isOnWhatsappTool := IsOnWhatsappTool(client)
	mcpServer.AddTool(isOnWhatsappTool, HandleIsOnWhatsapp(client))

	// Register get_chat_history tool
	getChatHistoryTool := GetChatHistoryTool(client)
	mcpServer.AddTool(getChatHistoryTool, HandleGetChatHistory(client))

	log.Println("Successfully registered 5 WhatsApp MCP tools:")
	log.Println("  - is_logged_in: Check authentication status")
	log.Println("  - get_qr_code: Generate QR code for login")
	log.Println("  - send_message: Send text messages")
	log.Println("  - is_on_whatsapp: Check phone number registration")
	log.Println("  - get_chat_history: Retrieve chat message history")
}
