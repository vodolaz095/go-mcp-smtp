package sender

import (
	"bytes"
	"context"
	"net/mail"
)

// Interface is interface Client satisfies
type Interface interface {
	Ping(ctx context.Context) error
	SendRaw(ctx context.Context, tos []*mail.Address, body *bytes.Buffer) error
	MakeBody(tos []*mail.Address, subject, body string) *bytes.Buffer
}
