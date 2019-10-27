GO=$(shell which go)

.PHONY: deps test lint install
.DEFAULT_GOAL := install

deps:
	@curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.21.0
	@$(GO) get ./...

lint:
	@echo golangci-lint...
	@golangci-lint run --tests=false --enable-all -D gochecknoglobals -D wsl -D funlen -D godox
	@echo test go install...
	@test -z $(shell $(GO) install)

test:
	@$(GO) test -v -race ./...

install:
	@$(GO) install
