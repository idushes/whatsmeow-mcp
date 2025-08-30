# WhatsApp MCP Server

A Model Context Protocol (MCP) server for WhatsApp functionality using the whatsmeow Go library.

## Features

This MCP server provides WhatsApp functionality through standardized tools that can be used by MCP clients like Claude Desktop, Cline, or other compatible applications.

### Currently Implemented Tools (Fake/Simulation)

- **is_logged_in** - Check if user is authenticated
- **get_qr_code** - Generate QR code for WhatsApp Web login
- **send_message** - Send text message to chat or contact
- **is_on_whatsapp** - Check if phone numbers are registered on WhatsApp
- **get_chat_history** - Get chat message history

> **Note:** These tools are currently implemented as simulations for testing the MCP server functionality. Real WhatsApp integration will be added in future versions.

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

### Running in stdio mode (for MCP clients)

```bash
go run main.go stdio
```

### Running in SSE mode (HTTP server)

```bash
go run main.go
```

The server will start on `http://localhost:3000` with the SSE endpoint available at `/sse`.

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

## API Documentation

### Tool: is_logged_in

Check if the WhatsApp client is authenticated.

**Parameters:** None

**Returns:**
```json
{
  "logged_in": false,
  "success": true
}
```

### Tool: get_qr_code

Generate QR code for WhatsApp Web login.

**Parameters:** None

**Returns:**
```json
{
  "qr_code": "2@ABC123...",
  "code": "2@ABC123...",
  "timeout": 30,
  "success": true
}
```

### Tool: send_message

Send a text message to a chat or contact.

**Parameters:**
- `to` (string, required): Recipient JID
- `text` (string, required): Message text content
- `quoted_message_id` (string, optional): ID of message to quote/reply to

**Returns:**
```json
{
  "message_id": "msg_1234567890",
  "timestamp": 1234567890,
  "success": true,
  "to": "1234567890@s.whatsapp.net",
  "text": "Hello World!"
}
```

### Tool: is_on_whatsapp

Check if phone numbers are registered on WhatsApp.

**Parameters:**
- `phones` (array of strings, required): Phone numbers in international format

**Returns:**
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

### Tool: get_chat_history

Get message history for a specific chat.

**Parameters:**
- `chat` (string, required): Chat JID
- `count` (number, optional): Number of messages to retrieve (default: 50)
- `before_message_id` (string, optional): Get messages before this ID

**Returns:**
```json
{
  "messages": [
    {
      "id": "msg_001",
      "from": "1234567890@s.whatsapp.net",
      "text": "Hello!",
      "timestamp": 1234567890,
      "chat": "1234567890@s.whatsapp.net"
    }
  ],
  "has_more": false,
  "success": true,
  "chat": "1234567890@s.whatsapp.net",
  "count": 1
}
```

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
├── main.go              # Main server implementation
├── config.env           # Configuration file
├── go.mod              # Go module definition
├── go.sum              # Go module checksums
├── README.md           # This file
└── MCP-TOOLS-PLAN.md   # Implementation plan and tool documentation
```

### Future Plans

- Real WhatsApp integration using whatsmeow library
- Additional tools for media handling, group management, etc.
- Persistent session management
- Enhanced error handling and logging
- Configuration validation

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
