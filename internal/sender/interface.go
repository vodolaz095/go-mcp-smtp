package sender

import "context"

// Interface is interface Client satisfies
type Interface interface {
	Ping(ctx context.Context) error
	SendRaw(ctx context.Context, recipients, subject, body string) error
}
