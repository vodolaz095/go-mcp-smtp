package sender

import (
	"context"
	"testing"
)

type Mock struct {
	T            *testing.T
	PingErr      error
	SendRawError error
}

func (m *Mock) Ping(context.Context) error {
	m.T.Helper()
	m.T.Logf("Ping is called")
	return m.PingErr
}

func (m *Mock) SendRaw(_ context.Context, recipients, subject, body string) error {
	m.T.Helper()
	m.T.Logf("SendRaw is called: recipients=%q, subject=%q, body=%q", recipients, subject, body)
	return m.SendRawError
}
