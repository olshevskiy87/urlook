GO=$(shell which go)

.PHONY: deps test lint install
.DEFAULT_GOAL := install

deps:
	@$(GO) get github.com/golangci/golangci-lint/cmd/golangci-lint
	@$(GO) get ./...

lint:
	@echo golangci-lint...
	@golangci-lint run --tests=false --enable-all -D gochecknoglobals ./...
	@echo test go install...
	@test -z $(shell $(GO) install)

test:
	@$(GO) test -v -race ./...

install:
	@$(GO) install
