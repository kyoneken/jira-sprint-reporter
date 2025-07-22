.PHONY: build build-all clean test fmt lint install help

# Build for current platform
build:
	go build -ldflags="-s -w" -o jira-sprint-reporter

# Build for all platforms
build-all:
	@echo "Building for all platforms..."
	@mkdir -p bin
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/jira-sprint-reporter-linux-amd64
	GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o bin/jira-sprint-reporter-linux-arm64
	GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o bin/jira-sprint-reporter-darwin-amd64
	GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o bin/jira-sprint-reporter-darwin-arm64
	GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o bin/jira-sprint-reporter-windows-amd64.exe
	@echo "Binaries built in bin/ directory"

# Clean build artifacts
clean:
	rm -f jira-sprint-reporter
	rm -rf bin/

# Run tests
test:
	go test -v ./...

# Format code
fmt:
	go fmt ./...

# Install dependencies
install:
	go mod tidy

# Run linter (requires golangci-lint to be installed)
lint:
	@which golangci-lint > /dev/null || (echo "golangci-lint not found. Install it with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest" && exit 1)
	golangci-lint run

# Show help
help:
	@echo "Available commands:"
	@echo "  build      - Build for current platform"
	@echo "  build-all  - Build for all platforms"
	@echo "  clean      - Clean build artifacts"
	@echo "  test       - Run tests"
	@echo "  fmt        - Format code"
	@echo "  lint       - Run linter"
	@echo "  install    - Install dependencies"
	@echo "  help       - Show this help"