package commands

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type SendRawInput struct {
	Recipients string `json:"name" jsonschema:"list of recipients in format of \"John Dow <john.dow@example.com>,Jane Dow <jane.dow@example.com>\""`
	Subject    string `json:"subject" jsonschema:"subject of email message, please use 8bit ANSI encoding"`
	Body       string `json:"body" jsonschema:"message body in plain text without any formatting"`
}

func (srv *MCP) SendRawEmail(ctx context.Context, _ *mcp.CallToolRequest, input SendRawInput) (*mcp.CallToolResult, Output, error) {
	err := srv.Sender.SendRaw(ctx, input.Recipients, input.Subject, input.Body)
	if err != nil {
		return nil, Output{Message: fmt.Sprintf("error sending message: %s", err)}, err
	}
	return nil, Output{Message: "message is accepted by submission server"}, nil
}
