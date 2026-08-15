package sender

import (
	"bytes"
	"fmt"
	"net/mail"
	"strings"
	"time"
)

// MakeBody renders plain text body for message
func (c *Client) MakeBody(tos []*mail.Address, subject, body string) *bytes.Buffer {
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
	return buh
}
