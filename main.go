package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"whatsmeow-mcp/internal/client"
	"whatsmeow-mcp/tools"

	"github.com/joho/godotenv"
	"github.com/mark3labs/mcp-go/server"
)

// Config holds the server configuration
type Config struct {
	Port          int
	Host          string
	ServerName    string
	ServerVersion string
	LogLevel      string
}

// loadConfig loads configuration from environment variables
func loadConfig() *Config {
	// Try to load .env file, ignore error if it doesn't exist
	_ = godotenv.Load(".env")

	config := &Config{
		Port:          3000,
		Host:          "localhost",
		ServerName:    "whatsmeow-mcp",
		ServerVersion: "1.0.0",
		LogLevel:      "info",
	}

	if port := os.Getenv("PORT"); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			config.Port = p
		}
	}

	if host := os.Getenv("HOST"); host != "" {
		config.Host = host
	}

	if name := os.Getenv("SERVER_NAME"); name != "" {
		config.ServerName = name
	}

	if version := os.Getenv("SERVER_VERSION"); version != "" {
		config.ServerVersion = version
	}

	if level := os.Getenv("LOG_LEVEL"); level != "" {
		config.LogLevel = level
	}

	return config
}

func main() {
	config := loadConfig()

	log.Printf("Starting %s v%s - WhatsApp MCP Server", config.ServerName, config.ServerVersion)
	log.Printf("Configuration: Host=%s, Port=%d, LogLevel=%s", config.Host, config.Port, config.LogLevel)

	// Initialize WhatsApp client
	whatsappClient := client.NewWhatsAppClient()
	log.Println("WhatsApp client initialized successfully")

	// Create MCP server with enhanced description
	mcpServer := server.NewMCPServer(
		config.ServerName,
		config.ServerVersion,
		server.WithInstructions("WhatsApp MCP Server - Provides WhatsApp functionality through standardized MCP tools. Enables AI agents and applications to send messages, check authentication status, verify phone numbers, retrieve chat history, and manage WhatsApp Web login via QR codes."),
	)

	// Register all WhatsApp tools
	tools.RegisterAllTools(mcpServer, whatsappClient)

	// Check if running in stdio mode (for MCP clients like Claude Desktop, Cline)
	if len(os.Args) > 1 && os.Args[1] == "stdio" {
		log.Println("Starting MCP server in stdio mode for direct client communication")
		log.Println("This mode is used by MCP clients like Claude Desktop and Cline")
		err := server.ServeStdio(mcpServer)
		if err != nil {
			log.Fatal("Failed to start stdio server:", err)
		}
	} else {
		// Run in SSE mode for HTTP-based communication
		log.Printf("Starting MCP server in SSE (Server-Sent Events) mode on %s:%d", config.Host, config.Port)
		log.Printf("SSE endpoint will be available at: http://%s:%d/sse", config.Host, config.Port)
		log.Println("This mode allows HTTP-based communication with the MCP server")

		sseServer := server.NewSSEServer(mcpServer,
			server.WithSSEEndpoint("/sse"),
		)
		err := sseServer.Start(fmt.Sprintf("%s:%d", config.Host, config.Port))
		if err != nil {
			log.Fatal("Failed to start SSE server:", err)
		}
	}
}
