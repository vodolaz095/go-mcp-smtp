package sender

import (
	"context"
	"fmt"
	"net/mail"
)

// SendRaw sends a raw email message through the SMTP server
func (c *Client) SendRaw(ctx context.Context, recipients, subject, body string) error {
	tos, err := mail.ParseAddressList(recipients)
	if err != nil {
		return fmt.Errorf("error parsing recipients %s : %w", recipients, err)
	}
	if len(tos) == 0 {
		return fmt.Errorf("empty list of recipients")
	}
	client, err := c.makeConnection(ctx)
	if err != nil {
		return fmt.Errorf("error making connection to %s %s as %s: %w", c.Network, c.Address, c.From, err)
	}
	defer client.Close()

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
