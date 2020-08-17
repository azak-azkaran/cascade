VERSION := $(shell git describe --always --long --dirty)
all: install

fetch:
	@go get -u github.com/stretchr/testify
	@go get -u ./...

build:
	@echo Building to current folder
	go build -v -ldflags="-X main.version=${VERSION}"

build_windows:
	@echo Building for Windows to current folder
	env GOOS=windows GOARCH=amd64 go build -v -ldflags="-X main.version=${VERSION}" -o cascade.exe

install: build
	@echo Installing to ${GOPATH}/bin
	go install

test: fetch
	@echo Running tests
	go test
	go test github.com/azak-azkaran/cascade/utils

coverage: fetch
	@echo Running Test with Coverage export
	go test github.com/azak-azkaran/cascade github.com/azak-azkaran/cascade/utils -v -coverprofile=cover.out -covermode=count
	go test github.com/azak-azkaran/cascade github.com/azak-azkaran/cascade/utils -json > report.json

coverall: coverage
	@echo Running Test with Coverall
	goveralls -coverprofile=cover.out -service=travis-ci -repotoken ${COVERALLS_TOKEN }

daemon: build
	@echo Moving cascade to /usr/local/bin
	mv ./cascade /usr/local/bin/
	@echo Copying config to /etc/systemd/system/cascade.service
	cp ./cascade.service /etc/systemd/system/cascade.service
	@echo restarting systemd
	systemctl daemon-reload
	@echo starting cascade as daemon
	systemctl start cascade

clean:
	go clean
