package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"time"
	"whatsmeow-mcp/internal/client"
	"whatsmeow-mcp/internal/database"
	"whatsmeow-mcp/internal/qrcode"
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

	// Database configuration
	DatabaseURL string
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

		// Database default
		DatabaseURL: "postgres://postgres:postgres@localhost:5432/whatsmeow_mcp?sslmode=disable",
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

	// Database configuration
	if databaseURL := os.Getenv("DATABASE_URL"); databaseURL != "" {
		config.DatabaseURL = databaseURL
	}

	return config
}

// maskDatabaseURL masks the password in database URL for logging
func maskDatabaseURL(databaseURL string) string {
	parsedURL, err := url.Parse(databaseURL)
	if err != nil {
		return "invalid-url"
	}

	if parsedURL.User != nil {
		username := parsedURL.User.Username()
		parsedURL.User = url.User(username) // Remove password
	}

	return parsedURL.String()
}

func main() {
	config := loadConfig()

	log.Printf("Starting %s v%s - WhatsApp MCP Server", config.ServerName, config.ServerVersion)
	log.Printf("Configuration: Host=%s, Port=%d, LogLevel=%s", config.Host, config.Port, config.LogLevel)
	log.Printf("Database URL: %s", maskDatabaseURL(config.DatabaseURL))

	// Create database if it doesn't exist
	if err := database.CreateDatabase(config.DatabaseURL); err != nil {
		log.Printf("Warning: Failed to create database: %v", err)
	}

	// Connect to database
	db, err := database.Connect(config.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize QR code generator
	staticDir := filepath.Join(".", "static")
	baseURL := fmt.Sprintf("http://%s:%d", config.Host, config.Port+1)
	qrGenerator := qrcode.NewQRCodeGenerator(staticDir, baseURL)

	// Start periodic cleanup of expired QR codes (every 10 minutes, remove files older than 30 minutes)
	qrGenerator.StartPeriodicCleanup(10*time.Minute, 30*time.Minute)

	// Initialize WhatsApp client (this will create whatsmeow tables)
	whatsappClient, err := client.NewWhatsmeowClient(db)
	if err != nil {
		log.Fatalf("Failed to initialize WhatsApp client: %v", err)
	}
	log.Println("WhatsApp client initialized successfully")

	// Run our custom migrations after whatsmeow has created its tables
	migrationsPath := filepath.Join(".", "migrations")
	if err := database.RunMigrations(db, migrationsPath); err != nil {
		log.Printf("Warning: Failed to run migrations: %v", err)
	}
	log.Println("Database connection established and migrations completed")

	// Create MCP server with enhanced description
	mcpServer := server.NewMCPServer(
		config.ServerName,
		config.ServerVersion,
		server.WithInstructions("WhatsApp MCP Server - Provides WhatsApp functionality through standardized MCP tools. Enables AI agents and applications to send messages, check authentication status, verify phone numbers, retrieve chat history, and manage WhatsApp Web login via QR codes."),
	)

	// Register all WhatsApp tools
	tools.RegisterAllTools(mcpServer, whatsappClient, qrGenerator)

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
		log.Printf("Static files endpoint will be available at: http://%s:%d/static/", config.Host, config.Port)
		log.Println("This mode allows HTTP-based communication with the MCP server")

		// Set up HTTP server with static file serving
		mux := http.NewServeMux()

		// Serve static files (QR codes)
		mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(staticDir))))

		// Create and start SSE server with custom handler
		sseServer := server.NewSSEServer(mcpServer,
			server.WithSSEEndpoint("/sse"),
		)

		// Start server with custom mux that includes static file serving
		go func() {
			log.Fatal("Static file server error:", http.ListenAndServe(fmt.Sprintf("%s:%d", config.Host, config.Port+1), http.StripPrefix("/static/", http.FileServer(http.Dir(staticDir)))))
		}()

		log.Printf("Static files served on: http://%s:%d/", config.Host, config.Port+1)

		// Start SSE server
		err := sseServer.Start(fmt.Sprintf("%s:%d", config.Host, config.Port))
		if err != nil {
			log.Fatal("Failed to start SSE server:", err)
		}
	}
}
