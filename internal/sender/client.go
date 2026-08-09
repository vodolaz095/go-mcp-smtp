package sender

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/mail"
	"net/smtp"
	"strings"
	"time"
)

// Client represents an SMTP client configuration for sending emails
type Client struct {
	Network  string
	Address  string
	Helo     string
	Host     string
	Username string
	Password string
	StartTLS bool
	From     string
}

func (c *Client) makeBody(tos []*mail.Address, subject, body string) []byte {
	to := make([]string, len(tos))
	for i := range tos {
		to[i] = tos[i].String()
	}

	now := time.Now()
	buh := bytes.NewBuffer(nil)
	fmt.Fprintf(buh, "Date: %s\r\n", now.Format(time.RFC1123Z))
	fmt.Fprintf(buh, "From: %s\r\n", c.From)
	fmt.Fprintf(buh, "To: %s\r\n", strings.Join(to, ","))
	fmt.Fprintf(buh, "Subject: %s\r\n", subject)
	fmt.Fprintf(buh, "X-Mailer: github.com/vodolaz095/go-mcp-smtp\r\n")
	fmt.Fprintf(buh, "Content-Type: text/plain; charset=\"utf-8\"\r\n")
	fmt.Fprintf(buh, "MIME-Version: 1.0\r\n")
	fmt.Fprintf(buh, "Message-Id: <%s@%s>\r\n", now.Format("20060102150405"), c.Host)
	fmt.Fprint(buh, "\r\n")
	fmt.Fprint(buh, body)
	return buh.Bytes()
}

func (c *Client) makeConnection(ctx context.Context) (client *smtp.Client, err error) {
	var myDialer net.Dialer

	con, err := myDialer.DialContext(ctx, c.Network, c.Address)
	if err != nil {
		return nil, fmt.Errorf("error dialing %s %s: %w", c.Network, c.Address, err)
	}
	client, err = smtp.NewClient(con, c.Host)
	if err != nil {
		return nil, fmt.Errorf("error starting smtp client for %s %s: %w", c.Network, c.Address, err)
	}
	err = client.Hello(c.Helo)
	if err != nil {
		return nil, fmt.Errorf("error sending helo for %s %s: %w", c.Network, c.Address, err)
	}
	if c.StartTLS {
		err = client.StartTLS(&tls.Config{ServerName: c.Host})
		if err != nil {
			return nil, fmt.Errorf("error starting tls for %s %s: %w", c.Network, c.Address, err)
		}
	}
	if c.Username != "" && c.Password != "" {
		err = client.Auth(smtp.PlainAuth("", c.Username, c.Password, c.Host))
		if err != nil {
			return nil, fmt.Errorf("error authenticating for %s %s as %s: %w", c.Network, c.Address, c.Username, err)
		}
	}
	err = client.Mail(c.From)
	if err != nil {
		return nil, fmt.Errorf("error setting MAIL FROM for %s %s as %s: %w", c.Network, c.Address, c.From, err)
	}
	return client, nil
}
