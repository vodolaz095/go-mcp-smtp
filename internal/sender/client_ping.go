package sender

import (
	"context"
	"fmt"
)

func (c *Client) Ping(ctx context.Context) error {
	client, err := c.makeConnection(ctx)
	if err != nil {
		return fmt.Errorf("error making connection to %s %s as %s: %w", c.Network, c.Address, c.From, err)
	}
	defer client.Close()
	return client.Quit()
}
