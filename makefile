VERSION := $(shell git describe --always --long --dirty)
all: install

build:
	Building cascade
	go build -i -v -ldflags="-X main.version=${VERSION}" 
	@go build -i -v -ldflags="-X main.version=${VERSION}" 

install:
	Installing cascade
	go install -i -v -ldflags="-X main.version=${VERSION}" 
	@go install -i -v -ldflags="-X main.version=${VERSION}" 

test:
	Running tests
	@go list -f '{{if len .TestGoFiles}}"go test  {{.ImportPath}}"{{end}}' ./... | xargs -L 1 sh -c

clean:
	@go clean
