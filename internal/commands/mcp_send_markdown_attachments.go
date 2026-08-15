package commands

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"mime/multipart"
	"net/textproto"
	"os"
	"path/filepath"

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

// SendMarkDownWithAttachments sends a multipart email message through the SMTP submission server with content rendered from markdown and optional attachments
func (srv *MCP) SendMarkDownWithAttachments(ctx context.Context, _ *mcp.CallToolRequest, input SendMarkDownWithAttachmentsInput) (*mcp.CallToolResult, Output, error) {
	var buf bytes.Buffer

	plain := bytes.NewBufferString(input.Markdown)
	rendered := mdToHTML(plain.Bytes())

	writer := multipart.NewWriter(&buf)

	buf.WriteString(fmt.Sprintf("Content-Type: multipart/alternative; boundary=%s\r\n\r\n", writer.Boundary()))

	textHeader := make(textproto.MIMEHeader)
	textHeader.Set("Content-Type", "text/plain; charset=UTF-8")
	textPart, err := writer.CreatePart(textHeader)
	if err != nil {
		err = fmt.Errorf("error creating plain text stream: %w", err)
		return nil, Output{Message: err.Error()}, err
	}
	_, err = textPart.Write(plain.Bytes())
	if err != nil {
		err = fmt.Errorf("error writing plain text part: %w", err)
		return nil, Output{Message: err.Error()}, err
	}

	htmlHeader := make(textproto.MIMEHeader)
	htmlHeader.Set("Content-Type", "text/html; charset=UTF-8")
	htmlPart, err := writer.CreatePart(htmlHeader)
	if err != nil {
		err = fmt.Errorf("error creating html stream: %w", err)
		return nil, Output{Message: err.Error()}, err
	}
	_, err = htmlPart.Write(rendered)
	if err != nil {
		err = fmt.Errorf("error writing plain text part: %w", err)
		return nil, Output{Message: err.Error()}, err
	}

	for i := range input.Attachments {
		data, fileErr := os.OpenFile(input.Attachments[i], os.O_RDONLY, 0644)
		if fileErr != nil {
			fileErr = fmt.Errorf("error opening file: %w", fileErr)
			return nil, Output{Message: fileErr.Error()}, fileErr
		}
		fileName := filepath.Base(input.Attachments[i])
		attachmentHeader := make(textproto.MIMEHeader)
		attachmentHeader.Set("Content-Type", fmt.Sprintf("application/octet-stream; name=%s", fileName))
		attachmentHeader.Set("Content-Transfer-Encoding", "base64")
		attachmentHeader.Set("Content-Disposition", fmt.Sprintf("filename=%s", fileName))
		attachmentPart, attachErr := writer.CreatePart(attachmentHeader)
		if attachErr != nil {
			attachErr = fmt.Errorf("error creating attachment stream: %w", attachErr)
			return nil, Output{Message: attachErr.Error()}, attachErr
		}
		_, attachErr = io.Copy(base64.NewEncoder(base64.StdEncoding, attachmentPart), data)
		if attachErr != nil {
			attachErr = fmt.Errorf("error sending file to attachment stream: %w", attachErr)
			return nil, Output{Message: attachErr.Error()}, attachErr
		}
	}

	err = writer.Close()
	if err != nil {
		err = fmt.Errorf("error closing writer: %w", err)
		return nil, Output{Message: err.Error()}, err
	}

	err = srv.Sender.SendRaw(ctx, input.Recipients, input.Subject, buf.String(),
		fmt.Sprintf("Content-Type: multipart/mixed; boundary=%s", writer.Boundary()),
	)
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
