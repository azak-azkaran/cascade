VERSION := $(shell git describe --always --long --dirty)
all: install

fetch:
	@go get ./...

build: fetch
	@echo Building to current folder
	go build -i -v -ldflags="-X main.version=${VERSION}" 

install: build
	@echo Installing to ${GOPATH}/bin
	go install

test: fetch
	@echo Running tests
	go list -f '{{if len .TestGoFiles}}"go test  {{.ImportPath}}"{{end}}' ./... | xargs -L 1 sh -c

coverage: fetch
	@echo Running Test with Coverage export
	go test -coverprofile=cover.out
	go test -json > report.json
	go test github.com/azak-azkaran/cascade/utils -coverprofile=./utils/cover.out
	go test github.com/azak-azkaran/cascade/utils -json > ./utils/report.json
	cd ../

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
