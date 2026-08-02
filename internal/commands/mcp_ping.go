package commands

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// PingInput is the input structure for the Ping method, currently empty as no parameters are needed
type PingInput struct{}

// Ping checks the connectivity and availability of the SMTP submission server
func (srv *MCP) Ping(ctx context.Context, _ *mcp.CallToolRequest, _ PingInput) (*mcp.CallToolResult, Output, error) {
	err := srv.Sender.Ping(ctx)
	if err != nil {
		return nil, Output{Message: fmt.Sprintf("error checking smtp submission server: %s", err)}, err
	}
	return nil, Output{Message: "smtp submission server is cooperating"}, nil
}
