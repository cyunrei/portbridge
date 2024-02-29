PROJECT_NAME := portbridge
GO := go
VERSION := $(shell git describe --tags)
LDFLAGS := -ldflags "-X main.version=$(subst v,,$(VERSION))"
BUILD_DIR := build

default: build

.PHONY: build
build:
	@$(GO) mod download
	@CGO_ENABLED=0 $(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(PROJECT_NAME)

.PHONY: build-all
build-all: build-linux build-windows build-darwin

.PHONY: build-linux
build-linux:
	@CGO_ENABLED=0 GOOS=linux GOARCH=arm64 $(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(PROJECT_NAME)-linux-aarch64
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(PROJECT_NAME)-linux-x86_64

.PHONY: build-windows
build-windows:
	@CGO_ENABLED=0 GOOS=windows GOARCH=arm64 $(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(PROJECT_NAME)-windows-aarch64.exe
	@CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(PROJECT_NAME)-windows-x86_64.exe

.PHONY: build-darwin
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

.PHONY: archive
archive:
	cd $(BUILD_DIR) && tar -czvf $(PROJECT_NAME)-linux-aarch64.tar.gz $(PROJECT_NAME)-linux-aarch64
	cd $(BUILD_DIR) && tar -czvf $(PROJECT_NAME)-linux-x86_64.tar.gz $(PROJECT_NAME)-linux-x86_64
	cd $(BUILD_DIR) && tar -czvf $(PROJECT_NAME)-windows-aarch64.tar.gz $(PROJECT_NAME)-windows-aarch64.exe
	cd $(BUILD_DIR) && tar -czvf $(PROJECT_NAME)-windows-x86_64.tar.gz $(PROJECT_NAME)-windows-x86_64.exe
	cd $(BUILD_DIR) && tar -czvf $(PROJECT_NAME)-darwin-aarch64.tar.gz $(PROJECT_NAME)-darwin-aarch64
	cd $(BUILD_DIR) && tar -czvf $(PROJECT_NAME)-darwin-x86_64.tar.gz $(PROJECT_NAME)-darwin-x86_64

	cd $(BUILD_DIR) && sha256sum $(PROJECT_NAME)-linux-aarch64.tar.gz > $(PROJECT_NAME)-linux-aarch64.tar.gz.sha256
	cd $(BUILD_DIR) && sha256sum $(PROJECT_NAME)-linux-x86_64.tar.gz > $(PROJECT_NAME)-linux-x86_64.tar.gz.sha256
	cd $(BUILD_DIR) && sha256sum $(PROJECT_NAME)-windows-aarch64.tar.gz > $(PROJECT_NAME)-windows-aarch64.tar.gz.sha256
	cd $(BUILD_DIR) && sha256sum $(PROJECT_NAME)-windows-x86_64.tar.gz > $(PROJECT_NAME)-windows-x86_64.tar.gz.sha256
	cd $(BUILD_DIR) && sha256sum $(PROJECT_NAME)-darwin-aarch64.tar.gz > $(PROJECT_NAME)-darwin-aarch64.tar.gz.sha256
	cd $(BUILD_DIR) && sha256sum $(PROJECT_NAME)-darwin-x86_64.tar.gz > $(PROJECT_NAME)-darwin-x86_64.tar.gz.sha256

.PHONY: release
release: build-all archive
