BINARY_NAME=genie
BUILD_DIR=bin

ifeq ($(OS),Windows_NT)
	BINARY_OUT=$(BINARY_NAME).exe
	INSTALL_DIR?=$(USERPROFILE)\bin
	CP=copy
	MKDIR=mkdir
else
	BINARY_OUT=$(BINARY_NAME)
	INSTALL_DIR?=$(HOME)/.local/bin
	CP=cp
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

install: build test
	@echo "Installing to $(INSTALL_DIR)..."
	@$(MKDIR) $(INSTALL_DIR)
	$(CP) $(BUILD_DIR)/$(BINARY_OUT) $(INSTALL_DIR)/$(BINARY_OUT)
ifeq ($(OS),Windows_NT)
	@echo "Adding $(INSTALL_DIR) to PATH..."
	@setx PATH "%PATH%;$(INSTALL_DIR)"
	@echo "Installation complete! Please RESTART your terminal to use genie."
else
	@echo "Installation complete!"
endif

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
