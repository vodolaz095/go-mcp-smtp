package commands

import (
	"github.com/vodolaz095/go-mcp-smtp/internal/sender"
)

type MCP struct {
	Sender *sender.Client
}
