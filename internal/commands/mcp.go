package commands

import (
	"github.com/vodolaz095/go-mcp-smtp/internal/sender"
)

// MCP is the Model Context Protocol server for SMTP operations
type MCP struct {
	Sender sender.Interface
}

// Output is the standard response format for MCP server operations
type Output struct {
	Message string `json:"message" jsonschema:"response"`
}
