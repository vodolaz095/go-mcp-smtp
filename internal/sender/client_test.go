package sender

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient(t *testing.T) {
	env := []string{
		"SMTP_HOST",
		"SMTP_USERNAME",
		"SMTP_PASSWORD",
		"SMTP_FROM",
		"SMTP_TO",
	}
	for i := range env {
		if os.Getenv(env[i]) == "" {
			t.Skipf("Environment variable %s is empty", env[i])
			return
		}
	}
	client := Client{
		Network:  "tcp",
		Address:  os.Getenv("SMTP_HOST") + ":587",
		Helo:     "localhost",
		Host:     os.Getenv("SMTP_HOST"),
		Username: os.Getenv("SMTP_USERNAME"),
		Password: os.Getenv("SMTP_PASSWORD"),
		StartTLS: true,
		From:     os.Getenv("SMTP_FROM"),
	}
	t.Run("ping", func(tt *testing.T) {
		err := client.Ping(tt.Context())
		if err != nil {
			tt.Errorf("error pinging: %s", err)
		}
	})
	t.Run("sendRawEmpty", func(tt *testing.T) {
		err := client.SendRaw(tt.Context(), "", "Test email send via go-mcp-smtp", "Test email send via go-mcp-smtp")
		assert.NotNil(tt, err)
		tt.Logf("error: %v", err)
		assert.ErrorContains(tt, err, "error parsing recipients")
		assert.ErrorContains(tt, err, "mail: no address")
	})
	t.Run("sendRawMalformed", func(tt *testing.T) {
		err := client.SendRaw(tt.Context(), "not.an.email.address", "Test email send via go-mcp-smtp", "Test email send via go-mcp-smtp")
		assert.NotNil(tt, err)
		tt.Logf("error: %v", err)
		assert.ErrorContains(tt, err, "error parsing recipients not.an.email.address")
		assert.ErrorContains(tt, err, "mail: missing '@' or angle-addr")
	})
	t.Run("sendRawPartiallyMalformed", func(tt *testing.T) {
		err := client.SendRaw(tt.Context(), os.Getenv("SMTP_TO")+", not.an.email.address", "Test email send via go-mcp-smtp", "Test email send via go-mcp-smtp")
		assert.NotNil(tt, err)
		tt.Logf("error: %v", err)
		assert.ErrorContains(tt, err, "error parsing recipients")
		assert.ErrorContains(tt, err, "mail: missing '@' or angle-addr")
	})
	t.Run("sendRawOK", func(tt *testing.T) {
		err := client.SendRaw(tt.Context(), os.Getenv("SMTP_TO"), "Test email send via go-mcp-smtp", "Test email send via go-mcp-smtp")
		if err != nil {
			tt.Errorf("error sending test email: %s", err)
		}
	})
}
