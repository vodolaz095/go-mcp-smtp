export app=go-mcp-smtp
export majorVersion=1
export minorVersion=0

export arch=$(shell uname)-$(shell uname -m)
export gittip=$(shell git log --format='%h' -n 1)
export subver=$(shell hostname)_on_$(shell date -u '+%Y-%m-%d_%I:%M:%S%p')
export patchVersion=$(shell git log --format='%h' | wc -l)
export ver=$(majorVersion).$(minorVersion).$(patchVersion).$(gittip)-$(arch)

clean:
	rm -f build/$(app)
	rm -f build/$(app).sh
	rm -f build/*.deb
	rm -f build/*.rpm
	rm -f build/*.apk

deps:
	go mod download
	go mod verify
	go mod tidy

build: deps
	CGO_ENABLED=0 go build -ldflags "-w -s -X main.Subversion=$(subver) -X main.Version=$(ver)" -o build/$(app) main.go
	ls -lh build/


tools:
	which go
	which govulncheck
	which golint
	which staticcheck

# https://go.dev/blog/govulncheck
# install it by go install golang.org/x/vuln/cmd/govulncheck@latest
vuln:
	which govulncheck
	govulncheck ./...

lint:
	gofmt -w=true -s=true -l=true ./
#	golint ./...
	go vet ./...
	staticcheck ./...

check: vuln lint test

test:
	go test -v ./...

tparse:
#	go install github.com/mfridman/tparse@latest
	go test -cover -race ./... -json | tparse -all -smallscreen -progress

cover:
	go test -v --cover ./...


start: run

upd/pkg:
	go get -u github.com/vodolaz095/pkg@latest
	go mod tidy

pkg: clean build
	echo "Downloading bounce file..." # TODO - use limited selfhosted in order to fix license
	curl -v -o build/bounces.txt https://raw.githubusercontent.com/zone-eu/zone-mta/master/config/bounces.txt
#	cp contrib/bounces.txt build/bounces.txt
	echo "Generating bash completions..."
	./build/$(app) completion bash > build/$(app).sh
	echo "Generating manual..."
	mkdir -p build/man
	./build/$(app) man build/man/
	echo "Generating config examples..."
	./build/$(app) generate ./build/config_full.yaml

include make/*.mk

.PHONY: test build generate
