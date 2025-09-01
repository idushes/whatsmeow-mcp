package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"
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
	MCPPort       int // Port for MCP/SSE server
	RESTPort      int // Port for REST API (health checks, static files)
	Host          string
	ServerName    string
	ServerVersion string
	LogLevel      string

	// Database configuration
	DatabaseURL string
}

// HealthChecker holds components needed for health checks
type HealthChecker struct {
	db     *database.MessageStore
	client client.WhatsAppClientInterface
	ready  bool
}

// Global health checker instance
var healthChecker *HealthChecker

// setupHealthChecks sets up health check endpoints
func setupHealthChecks(mux *http.ServeMux) {
	// Liveness probe - checks if the application is running
	mux.HandleFunc("/health/live", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"status":"ok","timestamp":"%s"}`, time.Now().Format(time.RFC3339))
	})

	// Readiness probe - checks if the application is ready to serve traffic
	mux.HandleFunc("/health/ready", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if healthChecker == nil || !healthChecker.ready {
			w.WriteHeader(http.StatusServiceUnavailable)
			fmt.Fprintf(w, `{"status":"not ready","timestamp":"%s"}`, time.Now().Format(time.RFC3339))
			return
		}

		// Check database connectivity
		if err := database.HealthCheck(healthChecker.db.GetDB()); err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			fmt.Fprintf(w, `{"status":"database unhealthy","error":"%s","timestamp":"%s"}`,
				err.Error(), time.Now().Format(time.RFC3339))
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"status":"ready","timestamp":"%s"}`, time.Now().Format(time.RFC3339))
	})

	// Combined health endpoint for Docker health check
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"status":"healthy","timestamp":"%s"}`, time.Now().Format(time.RFC3339))
	})
}

// loadConfig loads configuration from environment variables
func loadConfig() *Config {
	// Try to load .env file, ignore error if it doesn't exist
	_ = godotenv.Load(".env")

	config := &Config{
		MCPPort:       3000,
		RESTPort:      3001,
		Host:          "localhost",
		ServerName:    "whatsmeow-mcp",
		ServerVersion: "1.0.0",
		LogLevel:      "info",

		// Database default
		DatabaseURL: "postgres://postgres:postgres@localhost:5432/whatsmeow_mcp?sslmode=disable",
	}

	// MCP_PORT - port for MCP/SSE server
	if mcpPort := os.Getenv("MCP_PORT"); mcpPort != "" {
		if p, err := strconv.Atoi(mcpPort); err == nil {
			config.MCPPort = p
		}
	}

	// REST_PORT - port for REST API (health checks, static files)
	if restPort := os.Getenv("REST_PORT"); restPort != "" {
		if p, err := strconv.Atoi(restPort); err == nil {
			config.RESTPort = p
		}
	}

	// PORT - legacy support, sets MCP_PORT if MCP_PORT is not set
	if port := os.Getenv("PORT"); port != "" && os.Getenv("MCP_PORT") == "" {
		if p, err := strconv.Atoi(port); err == nil {
			config.MCPPort = p
			config.RESTPort = p + 1
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
	log.Printf("Configuration: Host=%s, MCP_PORT=%d, REST_PORT=%d, LogLevel=%s", config.Host, config.MCPPort, config.RESTPort, config.LogLevel)
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
	baseURL := fmt.Sprintf("http://%s:%d", config.Host, config.RESTPort)
	qrGenerator := qrcode.NewQRCodeGenerator(staticDir, baseURL)

	// Note: QR code cleanup can be implemented as a separate goroutine if needed

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

	// Initialize health checker
	messageStore := database.NewMessageStore(db)
	healthChecker = &HealthChecker{
		db:     messageStore,
		client: whatsappClient,
		ready:  false,
	}

	// Create MCP server with enhanced description
	mcpServer := server.NewMCPServer(
		config.ServerName,
		config.ServerVersion,
		server.WithInstructions("WhatsApp MCP Server - Provides WhatsApp functionality through standardized MCP tools. Enables AI agents and applications to send messages, check authentication status, verify phone numbers, retrieve chat history, and manage WhatsApp Web login via QR codes."),
	)

	// Register all WhatsApp tools
	tools.RegisterAllTools(mcpServer, whatsappClient, qrGenerator)

	// Mark as ready after successful initialization
	healthChecker.ready = true
	log.Println("Application is ready to serve traffic")

	// Check if running in stdio mode (for MCP clients like Claude Desktop, Cline)
	if len(os.Args) > 1 && os.Args[1] == "stdio" {
		log.Println("Starting MCP server in stdio mode for direct client communication")
		log.Println("This mode is used by MCP clients like Claude Desktop and Cline")

		// Set up signal handling for graceful shutdown
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

		go func() {
			<-sigChan
			log.Println("Received shutdown signal, cleaning up...")
			healthChecker.ready = false
			os.Exit(0)
		}()

		err := server.ServeStdio(mcpServer)
		if err != nil {
			log.Printf("Failed to start stdio server: %v", err)
			os.Exit(1)
		}
	} else {
		// Run in SSE mode for HTTP-based communication
		log.Printf("Starting MCP server in SSE (Server-Sent Events) mode")
		log.Printf("MCP/SSE endpoint will be available at: http://%s:%d/sse", config.Host, config.MCPPort)
		log.Printf("Static files endpoint will be available at: http://%s:%d/static/", config.Host, config.RESTPort)
		log.Printf("Health endpoints available at: http://%s:%d/health/*", config.Host, config.RESTPort)
		log.Println("This mode allows HTTP-based communication with the MCP server")

		// Set up HTTP server with all endpoints on single port
		mux := http.NewServeMux()

		// Setup health check endpoints
		setupHealthChecks(mux)

		// Serve static files (QR codes)
		mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(staticDir))))

		// Create SSE server
		sseServer := server.NewSSEServer(mcpServer,
			server.WithSSEEndpoint("/sse"),
		)

		// Create separate HTTP server for health checks and static files
		restServer := &http.Server{
			Addr:    fmt.Sprintf("%s:%d", config.Host, config.RESTPort),
			Handler: mux,
		}

		// Start REST API server in background
		go func() {
			log.Printf("Starting REST API server (health + static) on %s", restServer.Addr)
			if err := restServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Printf("REST server error: %v", err)
			}
		}()

		// Start MCP/SSE server in background
		go func() {
			log.Printf("Starting MCP/SSE server on %s:%d", config.Host, config.MCPPort)
			if err := sseServer.Start(fmt.Sprintf("%s:%d", config.Host, config.MCPPort)); err != nil {
				log.Printf("MCP/SSE server error: %v", err)
			}
		}()

		// Set up graceful shutdown
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

		// Wait for shutdown signal
		<-sigChan
		log.Println("Received shutdown signal, starting graceful shutdown...")

		// Mark as not ready
		healthChecker.ready = false

		// Create shutdown context with timeout
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		// Shutdown REST API server
		log.Println("Shutting down REST API server...")
		if err := restServer.Shutdown(ctx); err != nil {
			log.Printf("REST server shutdown error: %v", err)
		}

		// Note: SSE server doesn't have graceful shutdown method, it will be terminated

		// Note: QR code cleanup would be stopped here if implemented

		// Close database connection
		log.Println("Closing database connection...")
		db.Close()

		log.Println("Graceful shutdown completed")
	}
}
