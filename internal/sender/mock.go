package sender

import (
	"context"
	"testing"
)

// Mock is used as placeholder for Client in unit tests
type Mock struct {
	T            *testing.T
	PingErr      error
	SendRawError error

	PingCalled    bool
	SendRawCalled bool
}

// Ping emulates checking remote SMTP Submission server
func (m *Mock) Ping(context.Context) error {
	m.T.Helper()
	m.T.Logf("Ping is called")
	m.PingCalled = true
	return m.PingErr
}

// SendRaw emulates sending message via remote SMTP Submission server
func (m *Mock) SendRaw(_ context.Context, recipients, subject, body string) error {
	m.T.Helper()
	m.T.Logf("SendRaw is called: recipients=%q, subject=%q, body=%q", recipients, subject, body)
	m.SendRawCalled = true
	return m.SendRawError
}
