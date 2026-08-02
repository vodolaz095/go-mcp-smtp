package commands

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vodolaz095/go-mcp-smtp/internal/sender"
)

func TestMCP_Ping_OK(t *testing.T) {
	tr := sender.Mock{
		T:            t,
		PingErr:      nil,
		SendRawError: nil,
	}
	mcp := MCP{Sender: &tr}

	_, resp, err := mcp.Ping(t.Context(), nil, PingInput{})
	assert.Nil(t, err)
	assert.Equal(t, "smtp submission server is cooperating", resp.Message)
	assert.True(t, tr.PingCalled)
	assert.False(t, tr.SendRawCalled)
}

func TestMCP_Ping_Error(t *testing.T) {
	tr := sender.Mock{
		T:            t,
		PingErr:      fmt.Errorf("something is broken"),
		SendRawError: nil,
	}
	mcp := MCP{Sender: &tr}

	_, resp, err := mcp.Ping(t.Context(), nil, PingInput{})
	assert.ErrorContains(t, err, "something is broken")
	assert.Equal(t, "error checking smtp submission server: something is broken", resp.Message)
	assert.True(t, tr.PingCalled)
	assert.False(t, tr.SendRawCalled)
}

func TestMCP_SendRawEmail_OK(t *testing.T) {
	tr := sender.Mock{
		T:            t,
		PingErr:      nil,
		SendRawError: nil,
	}
	mcp := MCP{Sender: &tr}

	_, resp, err := mcp.SendRawEmail(t.Context(), nil, SendRawInput{
		Recipients: "somebody@example.org",
		Subject:    "test email, please ignore",
		Body:       "test email, please ignore",
	})
	assert.Nil(t, err)
	assert.Equal(t, "message is accepted by submission server", resp.Message)
	assert.False(t, tr.PingCalled)
	assert.True(t, tr.SendRawCalled)
}

func TestMCP_SendRawEmail_Error(t *testing.T) {
	tr := sender.Mock{
		T:            t,
		PingErr:      nil,
		SendRawError: fmt.Errorf("something is broken"),
	}
	mcp := MCP{Sender: &tr}

	_, resp, err := mcp.SendRawEmail(t.Context(), nil, SendRawInput{
		Recipients: "somebody@example.org",
		Subject:    "test email, please ignore",
		Body:       "test email, please ignore",
	})
	assert.Error(t, err)
	assert.ErrorContains(t, err, "something is broken")
	assert.Equal(t, "error sending message: something is broken", resp.Message)
	assert.False(t, tr.PingCalled)
	assert.True(t, tr.SendRawCalled)
}
