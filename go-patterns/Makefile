.PHONY: test test-verbose test-short test-cover bench lint vet fmt clean help

## test: Run all tests with gotestsum (pkgname format)
test:
	gotestsum --format=pkgname ./...

## test-verbose: Run all tests with verbose output and race detection
test-verbose:
	go test -v -race -count=1 ./...

## test-short: Run tests skipping long-running ones
test-short:
	go test -v -short -count=1 ./...

## test-cover: Run tests with coverage profile and HTML output
test-cover:
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

## bench: Run all benchmarks
bench:
	go test -bench=. -benchmem -run=^$$ ./...

## lint: Run golangci-lint
lint:
	golangci-lint run ./...

## vet: Run go vet
vet:
	go vet ./...

## fmt: Format code with gofumpt and fix imports
fmt:
	gofumpt -l -w .
	goimports -local github.com/netologist/go-patterns -w .

## clean: Remove test artifacts and coverage files
clean:
	rm -f coverage.out coverage.html

## help: Show this help message
help:
	@echo "Usage: make [target]"
	@sed -n 's/^##//p' $(MAKEFILE_LIST) | column -t -s ':' | sed -e 's/^/ /'
