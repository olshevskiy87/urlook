GO=$(shell which go)

.PHONY: deps test lint install
.DEFAULT_GOAL := install

deps:
	@curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOPATH)/bin v1.21.0
	@$(GO) get ./...

lint:
	@echo golangci-lint...
	@golangci-lint run
	@echo test go install...
	@test -z $(shell $(GO) install)

test:
	@$(GO) test -v -race ./...

install:
	@$(GO) install
