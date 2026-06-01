.PHONY: test test-verbose test-short test-cover bench lint vet fmt clean help

## test: Run all tests with gotestsum (pkgname format) in all submodules
test:
	$(MAKE) -C distributed test
	$(MAKE) -C go-patterns test

## test-verbose: Run all tests with verbose output and race detection in all submodules
test-verbose:
	$(MAKE) -C distributed test-verbose
	$(MAKE) -C go-patterns test-verbose

## test-short: Run tests skipping long-running ones in all submodules
test-short:
	$(MAKE) -C distributed test-short
	$(MAKE) -C go-patterns test-short

## test-cover: Run tests with coverage profile and HTML output in all submodules
test-cover:
	$(MAKE) -C distributed test-cover
	$(MAKE) -C go-patterns test-cover

## bench: Run all benchmarks in all submodules
bench:
	$(MAKE) -C distributed bench
	$(MAKE) -C go-patterns bench

## lint: Run golangci-lint in all submodules
lint:
	$(MAKE) -C distributed lint
	$(MAKE) -C go-patterns lint

## vet: Run go vet in all submodules
vet:
	$(MAKE) -C distributed vet
	$(MAKE) -C go-patterns vet

## fmt: Format code with gofumpt and fix imports in all submodules
fmt:
	$(MAKE) -C distributed fmt
	$(MAKE) -C go-patterns fmt

## clean: Remove test artifacts and coverage files in all submodules
clean:
	$(MAKE) -C distributed clean
	$(MAKE) -C go-patterns clean

## help: Show this help message
help:
	@echo "Usage: make [target]"
	@sed -n 's/^##//p' $(MAKEFILE_LIST) | column -t -s ':' | sed -e 's/^/ /'
