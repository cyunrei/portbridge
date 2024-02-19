PROJECT_NAME := portbridge
GO := go
LDFLAGS := -ldflags "-X main.version=`git describe --tags`"
BUILD_DIR := build

default: build

build:
	@$(GO) mod download
	@CGO_ENABLED=0 $(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(PROJECT_NAME)

build-all: build-linux build-windows build-darwin

build-linux:
	@CGO_ENABLED=0 GOOS=linux GOARCH=arm64 $(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(PROJECT_NAME)-linux-aarch64
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(PROJECT_NAME)-linux-x86_64

build-windows:
	@CGO_ENABLED=0 GOOS=windows GOARCH=arm64 $(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(PROJECT_NAME)-windows-aarch64.exe
	@CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(PROJECT_NAME)-windows-x86_64.exe

build-darwin:
	@CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 $(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(PROJECT_NAME)-darwin-aarch64
	@CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(PROJECT_NAME)-darwin-x86_64

upgrade:
	@$(GO) get -u -v
	@$(GO) mod download
	@$(GO) mod tidy
	@$(GO) mod verify

run:
	@./$(BUILD_DIR)/$(PROJECT_NAME)

clean:
	@$(GO) clean
	@$(GO) mod tidy
	@rm -rf $(BUILD_DIR)
