# WhatsApp MCP Server

A Model Context Protocol (MCP) server that provides WhatsApp functionality through standardized tools. This server enables AI agents and MCP clients to interact with WhatsApp services in a structured and reliable way.

The server implements the latest MCP protocol specification (2025-06-18) with Streamable HTTP transport, supporting both stdio and HTTP-based communication modes.

## Features

This MCP server exposes WhatsApp capabilities as standardized MCP tools that can be seamlessly integrated with AI agents, Claude Desktop, Cline, and other MCP-compatible applications. All tools provide comprehensive error handling, parameter validation, and detailed responses.

### Available MCP Tools

- **is_logged_in** - Check WhatsApp authentication status and session validity
- **get_qr_code** - Generate QR code for WhatsApp Web login with automatic expiration handling
- **send_message** - Send text messages to contacts or groups with optional message quoting/replies
- **is_on_whatsapp** - Verify WhatsApp registration status for phone numbers in bulk
- **get_chat_history** - Retrieve conversation history with pagination support

This server provides full WhatsApp functionality through the whatsmeow library integration.

## Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd whatsmeow-mcp
```

2. Install dependencies:
```bash
go mod download
```

3. Configure the server (optional):
```bash
cp config.env.example config.env
# Edit config.env with your preferred settings
```

## Configuration

The server can be configured using environment variables or a `config.env` file:

```env
# Server Configuration
PORT=3000
HOST=localhost

# MCP Configuration
SERVER_NAME=whatsmeow-mcp
SERVER_VERSION=1.0.0

# Logging
LOG_LEVEL=info
```

## Usage

The server supports two transport modes as defined in the MCP specification:

### 1. stdio mode (for MCP clients like Claude Desktop, Cline)

```bash
go run main.go stdio
```

In this mode, the server communicates through standard input/output streams, making it ideal for direct integration with MCP clients.

### 2. Streamable HTTP mode (HTTP server)

```bash
go run main.go
```

The server will start on `http://localhost:3000` with the MCP endpoint available at `/mcp`.

In this mode:
- The server accepts JSON-RPC messages via HTTP POST to `/mcp`
- Supports Server-Sent Events (SSE) for streaming responses
- Implements session management with `Mcp-Session-Id` header
- Validates protocol version via `MCP-Protocol-Version` header
- Provides separate health check endpoints on port 3001

### Building

```bash
go build -o whatsmeow-mcp main.go
```

## MCP Client Configuration

### Claude Desktop

Add to your `claude_desktop_config.json`:

```json
{
  "mcpServers": {
    "whatsmeow": {
      "command": "/path/to/whatsmeow-mcp",
      "args": ["stdio"]
    }
  }
}
```

### Cline (VSCode Extension)

Configure in Cline settings:
- **Server Command**: `/path/to/whatsmeow-mcp stdio`
- **Protocol**: `stdio`

## MCP Tools Documentation

### Tool: is_logged_in

**Purpose:** Check WhatsApp authentication status and session validity  
**Use Case:** Verify if the client is authenticated before performing operations that require login  
**Parameters:** None

**Response:**
```json
{
  "logged_in": false,
  "success": true
}
```

**AI Agent Notes:** Always check login status before attempting to send messages. If not logged in, guide user to use get_qr_code tool first.

---

### Tool: get_qr_code

**Purpose:** Generate QR code for WhatsApp Web authentication  
**Use Case:** Initial authentication setup - user scans QR code with mobile WhatsApp app  
**Parameters:** None  
**Important:** QR codes expire after 30 seconds and auto-refresh

**Response:**
```json
{
  "qr_code": "2@ABC123DEF456...",
  "code": "2@ABC123DEF456...", 
  "image_url": "http://localhost:6679/static/qr_1234567890_abcd1234.png",
  "timeout": 30,
  "success": true,
  "expires_at": 1234567920
}
```

**AI Agent Notes:** Use the `image_url` to display QR code image directly to users. The image is automatically generated and hosted by the server. QR codes expire after 30 seconds (check `expires_at` timestamp). Image files are automatically cleaned up after 5 minutes.

---

### Tool: send_message

**Purpose:** Send text messages to WhatsApp contacts or groups  
**Use Case:** Core messaging functionality with optional reply/quote capability  
**Authentication:** Requires active login session

**Parameters:**
- `to` (string, required): WhatsApp JID (phone number with @s.whatsapp.net suffix)
- `text` (string, required): Message content to send
- `quoted_message_id` (string, optional): ID of message to reply to/quote

**Response:**
```json
{
  "message_id": "msg_1234567890",
  "timestamp": 1234567890,
  "success": true,
  "to": "1234567890@s.whatsapp.net",
  "text": "Hello World!",
  "quoted_message_id": "msg_123"
}
```

**AI Agent Notes:** Validate phone number format. Check authentication first. Use quoted_message_id for contextual replies.

---

### Tool: is_on_whatsapp

**Purpose:** Verify WhatsApp registration status for phone numbers  
**Use Case:** Bulk verification before sending messages, contact validation  
**Efficiency:** Supports multiple phone numbers in single request

**Parameters:**
- `phones` (array of strings, required): Phone numbers in international format (e.g., +1234567890)

**Response:**
```json
{
  "results": [
    {
      "phone": "+1234567890",
      "is_on_whatsapp": true,
      "jid": "+1234567890@s.whatsapp.net"
    }
  ],
  "success": true
}
```

**AI Agent Notes:** Use returned JID for send_message operations. Batch multiple phone numbers for efficiency.

---

### Tool: get_chat_history

**Purpose:** Retrieve conversation history with pagination support  
**Use Case:** Context gathering, conversation analysis, message search  
**Pagination:** Supports count limits and before_message_id for scrolling

**Parameters:**
- `chat` (string, required): WhatsApp JID of the conversation
- `count` (number, optional): Messages to retrieve (default: 50, max: 100)
- `before_message_id` (string, optional): Message ID for pagination (get messages before this point)

**Response:**
```json
{
  "messages": [
    {
      "id": "msg_001",
      "from": "1234567890@s.whatsapp.net",
      "to": "self",
      "text": "Hello!",
      "timestamp": 1234567890,
      "chat": "1234567890@s.whatsapp.net",
      "quoted_message_id": "msg_000"
    }
  ],
  "has_more": false,
  "success": true,
  "chat": "1234567890@s.whatsapp.net",
  "count": 1
}
```

**AI Agent Notes:** Use has_more field to determine if additional messages exist. Implement pagination with before_message_id for large conversations.

## Error Handling

All tools return standardized error responses:

```json
{
  "success": false,
  "error": {
    "code": "ERROR_CODE",
    "message": "Human readable error message",
    "details": "Additional error details if available"
  }
}
```

Common error codes:
- `NOT_LOGGED_IN`: Client is not authenticated
- `INVALID_PARAMETERS`: Invalid or missing parameters
- `NETWORK_ERROR`: Network connectivity issue

## Development

### Project Structure

```
whatsmeow-mcp/
├── main.go                     # Main server entry point and configuration
├── internal/
│   ├── types/
│   │   ├── params.go          # Tool parameter definitions
│   │   └── responses.go       # Response type definitions
│   └── client/
│       ├── interface.go       # WhatsApp client interface
│       └── whatsmeow.go       # WhatsApp client implementation using whatsmeow
├── tools/
│   ├── is_logged_in.go        # Authentication status tool
│   ├── get_qr_code.go         # QR code generation tool
│   ├── send_message.go        # Message sending tool
│   ├── is_on_whatsapp.go      # Phone number verification tool
│   ├── get_chat_history.go    # Chat history retrieval tool
│   └── registry.go            # Tool registration and management
├── example.env                # Example environment configuration
├── go.mod                     # Go module definition
├── go.sum                     # Go module checksums
├── README.md                  # This documentation
└── MCP-TOOLS-PLAN.md         # Implementation roadmap
```

### Future Plans

- Additional tools for media handling, group management, etc.
- Persistent session management
- Enhanced error handling and logging
- Configuration validation
- Advanced message formatting options

## License

MIT License - see LICENSE file for details.

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## Support

For issues and questions:
- Create an issue in the GitHub repository
- Check the MCP-TOOLS-PLAN.md for implementation status
- Review the logs for troubleshooting

