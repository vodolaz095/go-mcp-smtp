package main

import (
	"flag"
	"log"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/vodolaz095/pkg/stopper"

	"github.com/vodolaz095/go-mcp-smtp/internal/commands"
	"github.com/vodolaz095/go-mcp-smtp/internal/sender"
)

var (
	network, address, helo, host, username, password string
	startTLS                                         bool
	from                                             string

	Version    = "development"
	Subversion = "development"
)

func main() {
	ctx, cancel := stopper.New()
	defer cancel()

	flag.StringVar(&network, "network", "tcp", "how to dial smtp server, can be tcp,tcp4,tcp6,unix")
	flag.StringVar(&address, "address", "localhost:587", "smtp submission server connection string")
	flag.StringVar(&host, "host", "localhost", "tls server name to use in TLS negotiation")
	flag.StringVar(&helo, "helo", "localhost", "smtp submission server connection string, protocol can be smtp or smtps")
	flag.StringVar(&username, "username", "", "username to use for authentication - can be blank")
	flag.StringVar(&password, "password", "", "password to use for authentication - can be blank")
	flag.StringVar(&from, "from", "mcp@localhost", "FROM header for email messages")
	flag.BoolVar(&startTLS, "start_tls", true, "start tls if possible")
	flag.Parse()

	log.Printf("Starting go-mcp-smtp. Version: %s. Subversion: %s. Please, report bugs here: https://github.com/vodolaz095/go-mcp-smtp/issues",
		Version, Subversion,
	)

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

	err := transport.Ping(ctx)
	if err != nil {
		log.Fatalf("error checking connection for smtp server: %s", err)
		return
	}
	log.Printf("SMTP server %s is cooperating...", address)
	srv := commands.MCP{Sender: &transport}

	server := mcp.NewServer(&mcp.Implementation{Name: "go-mcp-smtp",
		Title:       "go-mcp-smtp",
		Description: "MCP server to send email messages via SMTP submission server",
		Version:     Version,
		WebsiteURL:  "https://github.com/vodolaz095/go-mcp-smtp",
	}, nil)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "sendRawEmail",
		Title:       "send raw email message - agent should provide list of recipients, subject and raw message text",
		Description: "send raw email message",
	}, srv.SendRawEmail)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "ping",
		Title:       "ping SMTP Submission server to ensure it works with parameters provided",
		Description: "ensure smtp submission server is functional",
	}, srv.Ping)

	err = server.Run(ctx, &mcp.StdioTransport{})
	if err != nil {
		log.Fatal(err)
	}
}
