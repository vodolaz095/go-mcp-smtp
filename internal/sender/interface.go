package sender

import "context"

type Interface interface {
	Ping(ctx context.Context) error
	SendRaw(ctx context.Context, recipients, subject, body string) error
}
