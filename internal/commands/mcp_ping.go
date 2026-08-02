package commands

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type PingInput struct{}

func (srv *MCP) Ping(ctx context.Context, _ *mcp.CallToolRequest, _ PingInput) (*mcp.CallToolResult, Output, error) {
	err := srv.Sender.Ping(ctx)
	if err != nil {
		return nil, Output{Message: fmt.Sprintf("error checking smtp submission server: %s", err)}, err
	}
	return nil, Output{Message: "smtp submission server is cooperating"}, nil
}
