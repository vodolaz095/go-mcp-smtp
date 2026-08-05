# go-mcp-smtp Agents

This document describes the MCP server agents provided by the `go-mcp-smtp` project.

## Overview

`go-mcp-smtp` is an MCP (Model Context Protocol) server that enables sending email messages via the SMTP Submission protocol. It provides agents that can be integrated into AI systems to send emails programmatically.

## Available Agents

### Ping

**Purpose**: Check connectivity and availability of the SMTP submission server.

**Input**: None

**Output**: 
- `message`: Status message indicating server availability

**Example Usage**:
```json
{
  "tool": "Ping",
  "input": {}
}
```

**Possible Responses**:
- Success: `{"message": "smtp submission server is cooperating"}`
- Error: `{"message": "error checking smtp submission server: <error details>"}`

### SendRawEmail

**Purpose**: Send a raw email message through the SMTP submission server.

**Input**:
- `recipients`: List of recipients in format "Name <email@domain.com>" (comma-separated)
- `subject`: Subject of the email message (8bit ANSI encoding)
- `body`: Message body in plain text without any formatting

**Output**:
- `message`: Status message indicating message acceptance

**Example Usage**:
```json
{
  "tool": "SendRawEmail",
  "input": {
    "recipients": "John Doe <john@example.com>, Jane Smith <jane@example.com>",
    "subject": "Test Message",
    "body": "This is a test email sent via go-mcp-smtp."
  }
}
```

**Possible Responses**:
- Success: `{"message": "message is accepted by submission server"}`
- Error: `{"message": "error sending message: <error details>"}`

## Configuration

The MCP server requires configuration via command-line arguments when starting:

- `--network`: Network type (tcp, unix, etc.)
- `--address`: SMTP server address with port (e.g., smtp.example.org:587)
- `--host`: SMTP server hostname
- `--helo`: HELO/EHLO hostname
- `--username`: SMTP authentication username
- `--password`: SMTP authentication password
- `--from`: Default sender email address
- `--start_tls`: Whether to use STARTTLS (true/false)
- `--verbose`: Use verbose mode

See the README.md for a complete configuration example.

## Available Skills

### lint

**Purpose**: Run linting on the Go codebase and automatically fix any issues found.

**Usage**:
1. Verify required tools are available (go, gopls, staticcheck)
2. Run `make lint` to identify issues
3. Analyze lint output and apply fixes
4. Verify fixes resolved all issues

**Dependencies**:
- Go compiler (go)
- Go language server (gopls)
- Staticcheck tool
- Make build system

### vuln-scan

**Purpose**: Scan the codebase for known vulnerabilities using `govulncheck`.

**Usage**:
1. Ensure required tools are available (`make tools`)
2. Ensure dependencies are downloaded (`make deps`)
3. Run vulnerability scan (`make vuln`)
4. Identify vulnerable modules from the report
5. Update vulnerable modules using `go get`
6. Re-run `make deps` and `make vuln` to verify resolution

**Dependencies**:
- `govulncheck` tool
- Go 1.26+
- Internet access to fetch vulnerability database

### run-tests

**Purpose**: Run unit tests for the project.

**Usage**:
- `make test`: Run all tests with verbose output
- `make cover`: Run tests with coverage reporting
- `make check`: Run linting before tests

**Prerequisites**:
- Required tools: `make tools`
- Dependencies: `make deps`

**Additional Options**:
- `make cover`: Generate test coverage report
- `make check`: Run linting before executing tests
