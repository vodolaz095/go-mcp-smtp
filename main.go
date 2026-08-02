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

	Version string
)

func main() {
	ctx, cancel := stopper.New()
	defer cancel()

	flag.StringVar(&network, "network", "tcp", "how to dial smtp server, can be tcp,tcp4,tcp6,unix")
	flag.StringVar(&address, "address", "localhost:587", "smtp submission server connection string")
	flag.StringVar(&helo, "helo", "localhost", "smtp submission server connection string, protocol can be smtp or smtps")
	flag.StringVar(&username, "username", "", "username to use for authentication - can be blank")
	flag.StringVar(&password, "username", "", "username to use for authentication - can be blank")
	flag.StringVar(&from, "from", "mcp@localhost", "FROM header for email messages")
	flag.BoolVar(&startTLS, "start_tls", true, "start tls if possible")
	flag.Parse()

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

	srv := commands.MCP{Sender: &transport}

	server := mcp.NewServer(&mcp.Implementation{Name: "go-mcp-smtp", Version: Version}, nil)
	mcp.AddTool(server, &mcp.Tool{Name: "sendRawEail", Description: "send raw email message"}, srv.SendRawEmail)
	if err := server.Run(ctx, &mcp.StdioTransport{}); err != nil {
		log.Fatal(err)
	}
}
