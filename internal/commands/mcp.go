package commands

import (
	"github.com/vodolaz095/go-mcp-smtp/internal/sender"
)

type MCP struct {
	Sender *sender.Client
}

type Output struct {
	Message string `json:"message" jsonschema:"response"`
}
