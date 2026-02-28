.PHONY: build clean test run fmt vet lint sec check all help install release snapshot

BINARY_NAME=asciinema2video
BUILD_DIR=bin
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")

all: check build

build:
	go build -ldflags "-s -w -X main.version=$(VERSION)" -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/asciinema2video

build-all: clean
	GOOS=linux GOARCH=amd64 go build -ldflags "-s -w -X main.version=$(VERSION)" -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 ./cmd/asciinema2video
	GOOS=linux GOARCH=arm64 go build -ldflags "-s -w -X main.version=$(VERSION)" -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 ./cmd/asciinema2video
	GOOS=darwin GOARCH=amd64 go build -ldflags "-s -w -X main.version=$(VERSION)" -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 ./cmd/asciinema2video
	GOOS=darwin GOARCH=arm64 go build -ldflags "-s -w -X main.version=$(VERSION)" -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 ./cmd/asciinema2video
	GOOS=freebsd GOARCH=amd64 go build -ldflags "-s -w -X main.version=$(VERSION)" -o $(BUILD_DIR)/$(BINARY_NAME)-freebsd-amd64 ./cmd/asciinema2video
	GOOS=windows GOARCH=amd64 go build -ldflags "-s -w -X main.version=$(VERSION)" -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe ./cmd/asciinema2video

clean:
	rm -rf $(BUILD_DIR) dist coverage.out coverage.html
	go clean

test:
	go test -v -race -cover ./internal/...

test-coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

fmt:
	go fmt ./...
	@echo "Code formatted"

vet:
	go vet ./...
	@echo "Vet passed"

lint:
	@which golangci-lint > /dev/null || (echo "Installing golangci-lint..." && go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
	golangci-lint run ./...
	@echo "Lint passed"

sec:
	@which gosec > /dev/null || (echo "Installing gosec..." && go install github.com/securego/gosec/v2/cmd/gosec@latest)
	gosec -quiet ./...
	@echo "Security check passed"

check: fmt vet

check-all: fmt vet lint sec

# Release targets
release:
	@which goreleaser > /dev/null || (echo "Installing goreleaser..." && go install github.com/goreleaser/goreleaser@latest)
	goreleaser release --clean

snapshot:
	@which goreleaser > /dev/null || (echo "Installing goreleaser..." && go install github.com/goreleaser/goreleaser@latest)
	goreleaser release --snapshot --clean

# Run targets
run: build
	./$(BUILD_DIR)/$(BINARY_NAME) -i demo.cast -o demo.mp4

run-gif: build
	./$(BUILD_DIR)/$(BINARY_NAME) -i demo.cast -o demo.gif

run-webp: build
	./$(BUILD_DIR)/$(BINARY_NAME) -i demo.cast -o demo.webp

run-webm: build
	./$(BUILD_DIR)/$(BINARY_NAME) -i demo.cast -o demo.webm

run-rounded: build
	./$(BUILD_DIR)/$(BINARY_NAME) -i demo.cast -o demo-rounded.mp4 --border --border-radius 12

run-transparent: build
	./$(BUILD_DIR)/$(BINARY_NAME) -i demo.cast -o demo-transparent.webm --border --border-radius 12 --transparent

install: build
	cp $(BUILD_DIR)/$(BINARY_NAME) $(GOPATH)/bin/

help:
	@echo "Available targets:"
	@echo ""
	@echo "Build:"
	@echo "  all        - Run checks and build"
	@echo "  build      - Build the binary"
	@echo "  build-all  - Build for all platforms"
	@echo "  clean      - Clean build artifacts"
	@echo "  install    - Install to GOPATH/bin"
	@echo ""
	@echo "Test:"
	@echo "  test       - Run tests with race detection"
	@echo "  test-coverage - Run tests with coverage report"
	@echo ""
	@echo "Quality:"
	@echo "  fmt        - Format code"
	@echo "  vet        - Run go vet"
	@echo "  lint       - Run golangci-lint"
	@echo "  sec        - Run security scan (gosec)"
	@echo "  check      - Run fmt and vet"
	@echo "  check-all  - Run fmt, vet, lint, and sec"
	@echo ""
	@echo "Release:"
	@echo "  release    - Create release with goreleaser"
	@echo "  snapshot   - Create snapshot release"
	@echo ""
	@echo "Run:"
	@echo "  run        - Build and run (mp4)"
	@echo "  run-gif    - Build and run (gif)"
	@echo "  run-webp   - Build and run (webp)"
	@echo "  run-webm   - Build and run (webm)"
	@echo "  run-rounded - Build and run with rounded corners"
	@echo "  run-transparent - Build and run with transparency"
