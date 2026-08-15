package sender

import (
	"bytes"
	"context"
	"fmt"
	"net/mail"
	"strings"
	"testing"
	"time"
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
func (m *Mock) SendRaw(_ context.Context, tos []*mail.Address, body *bytes.Buffer) error {
	m.T.Helper()
	m.T.Logf("SendRaw is called: recipients=%q , body=%q", tos, body.String())
	m.SendRawCalled = true
	return m.SendRawError
}

// MakeBody renders simple email body
func (m *Mock) MakeBody(tos []*mail.Address, subject, body string) *bytes.Buffer {
	to := make([]string, len(tos))
	for i := range tos {
		to[i] = tos[i].String()
	}

	now := time.Now()
	buh := bytes.NewBuffer(nil)
	fmt.Fprintf(buh, "Date: %s\r\n", now.Format(time.RFC1123Z))
	fmt.Fprintf(buh, "From: %s\r\n", "mock@localhost")
	fmt.Fprintf(buh, "To: %s\r\n", strings.Join(to, ","))
	fmt.Fprintf(buh, "Subject: %s\r\n", subject)
	fmt.Fprintf(buh, "X-Mailer: github.com/vodolaz095/go-mcp-smtp\r\n")
	fmt.Fprintf(buh, "Content-Type: text/plain; charset=\"utf-8\"\r\n")
	fmt.Fprintf(buh, "MIME-Version: 1.0\r\n")
	fmt.Fprintf(buh, "Message-Id: <%s@%s>\r\n", now.Format("20060102150405"), "localhost")
	fmt.Fprint(buh, "\r\n")
	fmt.Fprint(buh, body)
	return buh
}
