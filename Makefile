GO=$(shell which go)

.PHONY: deps test lint install
.DEFAULT_GOAL := install

deps:
	@$(GO) get github.com/golang/lint/golint
	@$(GO) get honnef.co/go/tools/cmd/megacheck
	@$(GO) get ./...

lint:
	@if [ ! -z "$(shell gofmt -s -l .)" ]; then echo "gofmt blames these files:" && gofmt -s -l . && exit 1; fi
	@$(GO) vet ./...
	@megacheck ./...
	@golint -set_exit_status ./...
	@test -z $(shell $(GO) install)

test:
	@$(GO) test -v

install:
	@$(GO) install
