BINARY_NAME=genie
BUILD_DIR=bin

ifeq ($(OS),Windows_NT)
	BINARY_OUT=$(BINARY_NAME).exe
	MKDIR=mkdir
else
	BINARY_OUT=$(BINARY_NAME)
	MKDIR=mkdir -p
endif

.PHONY: all build test clean run fmt help install

all: build

build:
	@echo "Building $(BINARY_NAME)..."
	@$(MKDIR) $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_OUT) cmd/genie/main.go

test: build
	@echo "Running tests..."
	go test ./internal/...

install: test
	@echo "Installing genie via go install..."
	go install ./cmd/genie
	@echo "Installation complete! Make sure $$(go env GOPATH)/bin is in your PATH."

clean:
	@echo "Cleaning up..."
	@rm -rf $(BUILD_DIR)
	@rm -f genie genie.exe

run: build
	./$(BUILD_DIR)/$(BINARY_NAME)

help:
	@echo "Available commands:"
	@echo "  make build   - Build the genie binary"
	@echo "  make test    - Run unit tests"
	@echo "  make clean   - Remove build artifacts"
	@echo "  make run     - Build and run genie"
	@echo "  make help    - Show this help message"
