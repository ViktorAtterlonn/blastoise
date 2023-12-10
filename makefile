# Simple Makefile for a Go project

# Build the application
all: build

build:
	@echo "Building..."
	@make arm64

arm64:
	@echo "Building for macos-arm64..."
	@GOOS=darwin GOARCH=arm64 go build -o main-macos-arm64 cmd/blastoise/main.go

linux:
	@echo "Building for linux..."
	@GOOS=linux GOARCH=amd64 go build -o main-linux cmd/blastoise/main.go

# Run the application
run:
	@go run cmd/blastoise/main.go

# Test the application
test:
	@echo "Testing..."
	@go test ./...

# Clean the binary
clean:
	@echo "Cleaning..."
	@rm -f main

# Live Reload
watch:
	@if [ -x "$(GOPATH)/bin/air" ]; then \
	    "$(GOPATH)/bin/air"; \
		@echo "Watching...";\
	else \
	    read -p "air is not installed. Do you want to install it now? (y/n) " choice; \
	    if [ "$$choice" = "y" ]; then \
			go install github.com/cosmtrek/air@latest; \
	        "$(GOPATH)/bin/air"; \
				@echo "Watching...";\
	    else \
	        echo "You chose not to install air. Exiting..."; \
	        exit 1; \
	    fi; \
	fi

.PHONY: all build run test clean