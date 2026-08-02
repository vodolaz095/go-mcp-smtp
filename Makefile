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


# go install golang.org/x/lint/golint@latest
# go install honnef.co/go/tools/cmd/staticcheck@latest
lint:
	gofmt -w=true -s=true -l=true ./
	golint ./...
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

run:
	go run main.go \
	  --network=tcp \
	  --address=localhost:587 \
	  --host=localhost \
	  --helo=localhost \
	  --username=username \
	  --password=password \
	  --from="mcp@localhost" \
	  --verbose=true \
	  --start_tls=true

start: run

upd/pkg:
	go get -u github.com/vodolaz095/pkg@latest
	go mod tidy

include make/*.mk

nfpm/rpm: build
	mkdir -p ~/rpmbuild/BUILD/
	mkdir -p ~/rpmbuild/BUILDROOT/
	mkdir -p ~/rpmbuild/RPMS/x86_64/
	mkdir -p ~/rpmbuild/SOURCES/
	mkdir -p ~/rpmbuild/SPECS/
	mkdir -p ~/rpmbuild/SRPMS/
	SEMVER=$(majorVersion).$(minorVersion).$(patchVersion) nfpm package -p rpm -t ~/rpmbuild/RPMS/x86_64/
	rpm --addsign ~/rpmbuild/RPMS/x86_64/$(app)-$(majorVersion).$(minorVersion).$(patchVersion)-1.x86_64.rpm
	rpm --checksig -v ~/rpmbuild/RPMS/x86_64/$(app)-$(majorVersion).$(minorVersion).$(patchVersion)-1.x86_64.rpm

nfpm/deb: build
	SEMVER=$(majorVersion).$(minorVersion).$(patchVersion) nfpm package -p deb -t ./build/

nfpm/apk: build
	SEMVER=$(majorVersion).$(minorVersion).$(patchVersion) nfpm package -p apk -t ./build/

nfpm/arch: build
	SEMVER=$(majorVersion).$(minorVersion).$(patchVersion) nfpm package -p archlinux -t ./build/

nfpm/ipk: build
	SEMVER=$(majorVersion).$(minorVersion).$(patchVersion) nfpm package -p ipk -t ./build/

nfpm/all: build
	mkdir -p ~/rpmbuild/BUILD/
	mkdir -p ~/rpmbuild/BUILDROOT/
	mkdir -p ~/rpmbuild/RPMS/x86_64/
	mkdir -p ~/rpmbuild/SOURCES/
	mkdir -p ~/rpmbuild/SPECS/
	mkdir -p ~/rpmbuild/SRPMS/
	SEMVER=$(majorVersion).$(minorVersion).$(patchVersion) nfpm package -p rpm -t ~/rpmbuild/RPMS/x86_64/
	rpm --addsign ~/rpmbuild/RPMS/x86_64/$(app)-$(majorVersion).$(minorVersion).$(patchVersion)-1.x86_64.rpm
	rpm --checksig -v ~/rpmbuild/RPMS/x86_64/$(app)-$(majorVersion).$(minorVersion).$(patchVersion)-1.x86_64.rpm
	SEMVER=$(majorVersion).$(minorVersion).$(patchVersion) nfpm package -p deb -t ./build/
	SEMVER=$(majorVersion).$(minorVersion).$(patchVersion) nfpm package -p apk -t ./build/
	SEMVER=$(majorVersion).$(minorVersion).$(patchVersion) nfpm package -p archlinux -t ./build/
	SEMVER=$(majorVersion).$(minorVersion).$(patchVersion) nfpm package -p ipk -t ./build/

.PHONY: test build generate
