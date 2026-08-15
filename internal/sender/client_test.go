package sender

import (
	"net/mail"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSenderClient(t *testing.T) {
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
		err := client.SendRaw(tt.Context(), []*mail.Address{},
			client.MakeBody([]*mail.Address{}, "Test email send via go-mcp-smtp", "Test email send via go-mcp-smtp"),
		)
		assert.NotNil(tt, err)
		tt.Logf("error: %v", err)
		assert.ErrorContains(tt, err, "empty list of recipients")
	})
	t.Run("sendRawMalformed", func(tt *testing.T) {
		err := client.SendRaw(tt.Context(),
			[]*mail.Address{
				{
					Address: "not.an.email.address",
				},
			},
			client.MakeBody([]*mail.Address{}, "Test email send via go-mcp-smtp", "Test email send via go-mcp-smtp"),
		)
		assert.NotNil(tt, err)
		tt.Logf("error: %v", err)
		assert.ErrorContains(tt, err, "Malformed e-mail address")
	})
	t.Run("sendRawPartiallyMalformed", func(tt *testing.T) {
		err := client.SendRaw(tt.Context(),
			[]*mail.Address{
				{
					Address: os.Getenv("SMTP_TO") + ", not.an.email.address",
				},
			},
			client.MakeBody([]*mail.Address{}, "Test email send via go-mcp-smtp", "Test email send via go-mcp-smtp"),
		)
		assert.NotNil(tt, err)
		tt.Logf("error: %v", err)
		assert.ErrorContains(tt, err, "Malformed e-mail address")
	})
	t.Run("sendRawOK", func(tt *testing.T) {
		err := client.SendRaw(tt.Context(),
			[]*mail.Address{
				{
					Address: os.Getenv("SMTP_TO"),
				},
			},
			client.MakeBody([]*mail.Address{
				{Name: "", Address: os.Getenv("SMTP_TO")},
			}, "Test email send via go-mcp-smtp", "Test email send via go-mcp-smtp"),
		)
		if err != nil {
			tt.Errorf("error sending test email: %s", err)
		}
	})
}
