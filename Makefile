GO   := GO15VENDOREXPERIMENT=1 go
pkgs  = $(shell $(GO) list ./... | grep -v /vendor/)

all: format build

test:
	@echo ">> running tests"
	@$(GO) test -short $(pkgs)

format:
	@echo ">> formatting code"
	@$(GO) fmt $(pkgs)

vet:
	@echo ">> vetting code"
	@$(GO) vet $(pkgs)

build:
	@echo ">> building binaries"
	@./scripts/build.sh

docker:
	@docker build -t goproject:$(shell git rev-parse --short HEAD) .

.PHONY: all format build test vet docker assets
