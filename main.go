package main

import (
	"context"
	"flag"
	"log"
	"os"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/vodolaz095/pkg/stopper"

	"github.com/vodolaz095/go-mcp-smtp/internal/commands"
	"github.com/vodolaz095/go-mcp-smtp/internal/sender"
)

var (
	network, address, helo, host, username, password string
	startTLS, verbose                                bool
	from                                             string

	// Version is the current version of the application, typically set during build
	Version = "development"
	// Subversion is the extra build information for the application, typically set during build
	Subversion = "development"
)

func main() {
	ctx, cancel := stopper.New()
	defer cancel()
	var err error

	flag.StringVar(&network, "network", "tcp", "how to dial smtp server, can be tcp,tcp4,tcp6,unix")
	flag.StringVar(&address, "address", "localhost:587", "smtp submission server connection string")
	flag.StringVar(&host, "host", "localhost", "tls server name to use in TLS negotiation")
	flag.StringVar(&helo, "helo", "localhost", "how to introduce ourselves during HELO step of SMTP negotiation")
	flag.StringVar(&username, "username", "", "username to use for authentication - can be blank")
	flag.StringVar(&password, "password", "", "password to use for authentication - can be blank")
	flag.StringVar(&from, "from", "mcp@localhost", "FROM header for email messages")
	flag.BoolVar(&startTLS, "start_tls", true, "start tls if possible")
	flag.BoolVar(&verbose, "verbose", false, "use verbose mode")
	flag.Parse()

	log.SetOutput(os.Stderr)
	if verbose {
		log.Printf("Starting go-mcp-smtp. Version: %s. Subversion: %s. Please, report bugs here: https://github.com/vodolaz095/go-mcp-smtp/issues",
			Version, Subversion,
		)
	}

	transport := sender.Client{
		Network:  network,
		Address:  address,
		Helo:     helo,
		Host:     host,
		Username: username,
		Password: password,
		StartTLS: startTLS,
		From:     from,
	}

	err = transport.Ping(ctx)
	if err != nil {
		log.Fatalf("error checking connection for smtp server: %s", err)
		return
	}
	if verbose {
		log.Printf("SMTP server %s is cooperating...", address)
	}
	srv := commands.MCP{Sender: &transport}

	server := mcp.NewServer(&mcp.Implementation{Name: "go-mcp-smtp",
		Title:       "go-mcp-smtp",
		Description: "MCP server to send email messages via SMTP submission server",
		Version:     "0.1.0",
		WebsiteURL:  "https://github.com/vodolaz095/go-mcp-smtp",
	}, &mcp.ServerOptions{
		InitializedHandler: func(ctx context.Context, _ *mcp.InitializedRequest) {
			if verbose {
				log.Println("MCP server is initialized!")
			}
		},
		Capabilities: &mcp.ServerCapabilities{
			Prompts: &mcp.PromptCapabilities{
				ListChanged: true,
			},
			Tools: &mcp.ToolCapabilities{
				ListChanged: true,
			},
			Resources: &mcp.ResourceCapabilities{
				ListChanged: true,
				Subscribe:   true,
			},
		},
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "sendRawEmail",
		Title:       "send raw email message - agent should provide list of recipients (in RFC 5322 format like `John Doe <john.doe@example.org>, Jane Doe <jane.doe@example.org>`), subject and raw message text",
		Description: "send raw email message",
	}, srv.SendRawEmail)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "ping",
		Title:       "ping SMTP Submission server to ensure it works with parameters provided",
		Description: "ensure smtp submission server is functional",
	}, srv.Ping)
	if verbose {
		log.Printf("Starting MCP...")
	}
	err = server.Run(ctx, &mcp.StdioTransport{})
	if err != nil {
		log.Fatal(err)
	}
}
