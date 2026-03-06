BINARY_NAME=genie
BUILD_DIR=bin

.PHONY: all build test clean run fmt help

all: build

build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) cmd/genie/main.go

test:
	@echo "Running tests..."
	go test ./internal/...

clean:
	@echo "Cleaning up..."
	@rm -rf $(BUILD_DIR)
	@rm -f genie

run: build
	./$(BUILD_DIR)/$(BINARY_NAME)

help:
	@echo "Available commands:"
	@echo "  make build   - Build the genie binary"
	@echo "  make test    - Run unit tests"
	@echo "  make clean   - Remove build artifacts"
	@echo "  make run     - Build and run genie"
	@echo "  make help    - Show this help message"
