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
	buh := bytes.NewBufferString("Date: " + now.Format(time.RFC1123Z) + "\r\n")
	buh.WriteString("From: " + c.From + "\r\n")
	buh.WriteString("To: " + strings.Join(to, ",") + "\r\n")
	buh.WriteString(fmt.Sprintf("Subject: %s\r\n", subject))
	buh.WriteString("X-Mailer: github.com/vodolaz095/go-mcp-smtp\r\n")
	buh.WriteString(fmt.Sprintf("Message-Id: <%s@localhost>\r\n", now.Format("20060102150405")))
	buh.WriteString("\r\n")
	buh.WriteString(body)
	return buh.Bytes()
}

func (c *Client) Send(ctx context.Context, recipients, subject, body string) error {
	tos, err := mail.ParseAddressList(recipients)
	if err != nil {
		return fmt.Errorf("error parsing recipients %s : %w", recipients, err)
	}
	if len(tos) == 0 {
		return fmt.Errorf("empty list of recipients")
	}

	var myDialer net.Dialer

	con, err := myDialer.DialContext(ctx, c.Network, c.Address)
	if err != nil {
		return fmt.Errorf("error dialing %s %s: %w", c.Network, c.Address, err)
	}
	client, err := smtp.NewClient(con, c.Host)
	if err != nil {
		return fmt.Errorf("error starting smtp client for %s %s: %w", c.Network, c.Address, err)
	}
	defer client.Close()
	err = client.Hello(c.Helo)
	if err != nil {
		return fmt.Errorf("error sending helo for %s %s: %w", c.Network, c.Address, err)
	}
	if c.StartTLS {
		err = client.StartTLS(&tls.Config{ServerName: c.Host})
		if err != nil {
			return fmt.Errorf("error starting tls for %s %s: %w", c.Network, c.Address, err)
		}
	}
	if c.Username != "" && c.Password != "" {
		err = client.Auth(smtp.PlainAuth("", c.Username, c.Password, c.Host))
		if err != nil {
			return fmt.Errorf("error authenticating for %s %s as %s: %w", c.Network, c.Address, c.Username, err)
		}
	}
	err = client.Mail(c.From)
	if err != nil {
		return fmt.Errorf("error setting MAIL FROM for %s %s as %s: %w", c.Network, c.Address, c.From, err)
	}
	for i := range tos {
		err = client.Rcpt(tos[i].Address)
		if err != nil {
			return fmt.Errorf("error setting RCPT TO for %s %s as %s: %w", c.Network, c.Address, c.From, err)
		}
	}
	wc, err := client.Data()
	if err != nil {
		return fmt.Errorf("error executing DATA for %s %s: %w", c.Network, c.Address, err)
	}
	_, err = wc.Write(c.makeBody(tos, subject, body))
	if err != nil {
		return fmt.Errorf("error sending email body for %s %s: %w", c.Network, c.Address, err)
	}

	err = wc.Close()
	if err != nil {
		return fmt.Errorf("error closing email body for %s %s: %w", c.Network, c.Address, err)
	}

	err = client.Quit()
	if err != nil {
		return fmt.Errorf("error quiting for %s %s: %w", c.Network, c.Address, err)
	}

	return nil
}
