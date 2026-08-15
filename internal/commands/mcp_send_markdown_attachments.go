package commands

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"mime/multipart"
	"net/mail"
	"net/textproto"
	"os"
	"path/filepath"
	"time"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// SendMarkDownWithAttachmentsInput contains the parameters for sending a markdown email message to list of recipients
type SendMarkDownWithAttachmentsInput struct {
	Recipients  string   `json:"recipients" jsonschema:"list of recipients in RFC 5322 format like \"John Dow <john.dow@example.com>, Jane Dow <jane.dow@example.com>\" or even \"jane.doe@example.org\" for single address."`
	Subject     string   `json:"subject" jsonschema:"subject of email message, please use 8bit ANSI encoding"`
	Markdown    string   `json:"markdown" jsonschema:"message body in github flavored markdown format"`
	Attachments []string `json:"attachments" jsonschema:"list of absolute file paths to attachment files, can be empty"`
}

// SendMarkDownWithAttachmentsEmail sends a multipart email message through the SMTP submission server with content rendered from markdown and optional attachments
func (srv *MCP) SendMarkDownWithAttachmentsEmail(ctx context.Context, _ *mcp.CallToolRequest, input SendMarkDownWithAttachmentsInput) (*mcp.CallToolResult, Output, error) {
	var buf bytes.Buffer
	now := time.Now()
	tos, err := mail.ParseAddressList(input.Recipients)
	if err != nil {
		err = fmt.Errorf("error parsing recipients %s : %w", input.Recipients, err)
		return nil, Output{Message: err.Error()}, err
	}
	if len(tos) == 0 {
		err = fmt.Errorf("empty list of recipients")
		return nil, Output{Message: err.Error()}, err
	}

	plain := bytes.NewBufferString(input.Markdown)
	rendered := mdToHTML(plain.Bytes())

	// Create the outer multipart/mixed writer
	writer := multipart.NewWriter(&buf)

	// Write headers
	buf.WriteString(fmt.Sprintf("Date: %s\r\n", now.Format(time.RFC1123Z)))
	buf.WriteString(fmt.Sprintf("From: %s\r\n", srv.From))
	buf.WriteString(fmt.Sprintf("To: %s\r\n", input.Recipients))
	buf.WriteString(fmt.Sprintf("Subject: %s\r\n", input.Subject))
	buf.WriteString(fmt.Sprintf("X-Mailer: github.com/vodolaz095/go-mcp-smtp\r\n"))
	buf.WriteString(fmt.Sprintf("MIME-Version: 1.0\r\n"))
	buf.WriteString(fmt.Sprintf("Message-Id: <%s@%s>\r\n", now.Format("20060102150405"), srv.Host))
	buf.WriteString(fmt.Sprintf("Content-Type: multipart/mixed; boundary=%s\r\n", writer.Boundary()))
	buf.WriteString(fmt.Sprintf("\r\n"))

	// Create the multipart/alternative part for text and HTML
	alternativeWriter := multipart.NewWriter(&buf)
	alternativeHeader := make(textproto.MIMEHeader)
	alternativeHeader.Set("Content-Type", fmt.Sprintf("multipart/alternative; boundary=%s", alternativeWriter.Boundary()))
	alternativePart, err := writer.CreatePart(alternativeHeader)
	if err != nil {
		err = fmt.Errorf("error creating multipart/alternative part: %w", err)
		return nil, Output{Message: err.Error()}, err
	}

	// Add text/plain part to the alternative container
	textHeader := make(textproto.MIMEHeader)
	textHeader.Set("Content-Type", "text/plain; charset=UTF-8")
	textPart, err := alternativeWriter.CreatePart(textHeader)
	if err != nil {
		err = fmt.Errorf("error creating plain text stream: %w", err)
		return nil, Output{Message: err.Error()}, err
	}
	_, err = textPart.Write(plain.Bytes())
	if err != nil {
		err = fmt.Errorf("error writing plain text part: %w", err)
		return nil, Output{Message: err.Error()}, err
	}

	// Add text/html part to the alternative container
	htmlHeader := make(textproto.MIMEHeader)
	htmlHeader.Set("Content-Type", "text/html; charset=UTF-8")
	htmlPart, err := alternativeWriter.CreatePart(htmlHeader)
	if err != nil {
		err = fmt.Errorf("error creating html stream: %w", err)
		return nil, Output{Message: err.Error()}, err
	}
	_, err = htmlPart.Write(rendered)
	if err != nil {
		err = fmt.Errorf("error writing html part: %w", err)
		return nil, Output{Message: err.Error()}, err
	}

	// Close the alternative writer and write it to the main buffer
	err = alternativeWriter.Close()
	if err != nil {
		err = fmt.Errorf("error closing alternative writer: %w", err)
		return nil, Output{Message: err.Error()}, err
	}
	_, err = alternativePart.Write(buf.Bytes()[buf.Len()-len(alternativeWriter.Boundary()):])
	if err != nil {
		err = fmt.Errorf("error writing alternative part: %w", err)
		return nil, Output{Message: err.Error()}, err
	}

	// Add attachments to the outer multipart/mixed container
	for i := range input.Attachments {
		data, fileErr := os.OpenFile(input.Attachments[i], os.O_RDONLY, 0644)
		if fileErr != nil {
			fileErr = fmt.Errorf("error opening file: %w", fileErr)
			data.Close()
			return nil, Output{Message: fileErr.Error()}, fileErr
		}

		fileName := filepath.Base(input.Attachments[i])
		attachmentHeader := make(textproto.MIMEHeader)
		attachmentHeader.Set("Content-Type", fmt.Sprintf("application/octet-stream; name=%s", fileName))
		attachmentHeader.Set("Content-Transfer-Encoding", "base64")
		attachmentHeader.Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
		attachmentPart, attachErr := writer.CreatePart(attachmentHeader)
		if attachErr != nil {
			attachErr = fmt.Errorf("error creating attachment stream: %w", attachErr)
			data.Close()
			return nil, Output{Message: attachErr.Error()}, attachErr
		}
		_, attachErr = io.Copy(base64.NewEncoder(base64.StdEncoding, attachmentPart), data)
		if attachErr != nil {
			attachErr = fmt.Errorf("error sending file to attachment stream: %w", attachErr)
			data.Close()
			return nil, Output{Message: attachErr.Error()}, attachErr
		}
		data.Close()
	}

	err = writer.Close()
	if err != nil {
		err = fmt.Errorf("error closing writer: %w", err)
		return nil, Output{Message: err.Error()}, err
	}

	err = srv.Sender.SendRaw(ctx, tos, &buf)
	if err != nil {
		return nil, Output{Message: fmt.Sprintf("error sending message: %s", err)}, err
	}
	return nil, Output{Message: "message is accepted by submission server"}, nil

}

func mdToHTML(md []byte) []byte {
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse(md)

	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	return markdown.Render(doc, renderer)
}
