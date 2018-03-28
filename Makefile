GO=$(shell which go)

.PHONY: all deps test install
.DEFAULT_GOAL := install

all: install

deps:
	@$(GO) get github.com/golang/lint/golint
	@$(GO) get honnef.co/go/tools/cmd/megacheck

test:
	@if [ ! -z "$(shell gofmt -s -l .)" ]; then echo "gofmt blames these files:" && gofmt -s -l . && exit 1; fi
	@go vet ./...
	@megacheck ./...
	@golint -set_exit_status ./...
	@test -z $(shell $(GO) install)

install:
	@$(GO) install
