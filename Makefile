export app=tagtaa
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
	rm -rf build/man/
	rm -rf build/examples/
	rm -f build/bounces.txt
	rm -f build/config.yaml
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

configs:
	test -f runtime.yaml || cp use_cases/development.yaml runtime.yaml

run: configs
	go run main.go start --config ./runtime.yaml

binary:
	build/$(app) --config ./runtime.yaml

start: run

blacklist: configs
	go run main.go extract --limit 0 --config ./runtime.yaml

check_interfaces:
	go run main.go check

upd/msmtpd:
	go get -u github.com/vodolaz095/msmtpd@latest
	go mod tidy

upd/pkg:
	go get -u github.com/vodolaz095/pkg@latest
	go mod tidy

upd/all: upd/msmtpd upd/pkg

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


nfpm/rpm: pkg
	mkdir -p ~/rpmbuild/BUILD/
	mkdir -p ~/rpmbuild/BUILDROOT/
	mkdir -p ~/rpmbuild/RPMS/x86_64/
	mkdir -p ~/rpmbuild/SOURCES/
	mkdir -p ~/rpmbuild/SPECS/
	mkdir -p ~/rpmbuild/SRPMS/
	SEMVER=$(majorVersion).$(minorVersion).$(patchVersion) nfpm package -p rpm -t ~/rpmbuild/RPMS/x86_64/
	rpm --addsign ~/rpmbuild/RPMS/x86_64/$(app)-$(majorVersion).$(minorVersion).$(patchVersion)-1.x86_64.rpm
	rpm --checksig -v ~/rpmbuild/RPMS/x86_64/$(app)-$(majorVersion).$(minorVersion).$(patchVersion)-1.x86_64.rpm

nfpm/deb: pkg
	SEMVER=$(majorVersion).$(minorVersion).$(patchVersion) nfpm package -p deb -t ./build/

nfpm/apk: pkg
	SEMVER=$(majorVersion).$(minorVersion).$(patchVersion) nfpm package -p apk -t ./build/

nfpm/arch: pkg
	SEMVER=$(majorVersion).$(minorVersion).$(patchVersion) nfpm package -p archlinux -t ./build/

nfpm/ipk: pkg
	SEMVER=$(majorVersion).$(minorVersion).$(patchVersion) nfpm package -p ipk -t ./build/

nfpm/all: pkg
	SEMVER=$(majorVersion).$(minorVersion).$(patchVersion) nfpm package -p rpm -t ./build/
	rpm --addsign ./build/$(app)-$(majorVersion).$(minorVersion).$(patchVersion)-1.x86_64.rpm
	rpm --checksig -v ./build/$(app)-$(majorVersion).$(minorVersion).$(patchVersion)-1.x86_64.rpm
	SEMVER=$(majorVersion).$(minorVersion).$(patchVersion) nfpm package -p deb -t ./build/
	SEMVER=$(majorVersion).$(minorVersion).$(patchVersion) nfpm package -p archlinux -t ./build/
	SEMVER=$(majorVersion).$(minorVersion).$(patchVersion) nfpm package -p apk -t ./build/
	SEMVER=$(majorVersion).$(minorVersion).$(patchVersion) nfpm package -p ipk -t ./build/


lmtp_holod:
	cd internal/sender/lmtp && go test -c -o lmtp_test . && scp lmtp_test holod.local:/home/vodolaz095/

include make/*.mk

.PHONY: test build generate
