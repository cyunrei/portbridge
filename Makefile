PROJECT_NAME := portbridge
GO := go
VERSION := $(shell git describe --tags)
LDFLAGS := -ldflags "-X main.version=$(subst v,,$(VERSION))"
BUILD_DIR := build
MAIN_FILE := cmd/main.go

default: build

.PHONY: build
build:
	@$(GO) mod download
	@CGO_ENABLED=0 $(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(PROJECT_NAME) $(MAIN_FILE)

.PHONY: build-all
build-all: build-linux build-windows build-darwin

.PHONY: build-linux
build-linux:
	@CGO_ENABLED=0 GOOS=linux GOARCH=arm64 $(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(PROJECT_NAME)-linux-aarch64 $(MAIN_FILE)
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(PROJECT_NAME)-linux-x86_64 $(MAIN_FILE)

.PHONY: build-windows
build-windows:
	@CGO_ENABLED=0 GOOS=windows GOARCH=arm64 $(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(PROJECT_NAME)-windows-aarch64.exe $(MAIN_FILE)
	@CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(PROJECT_NAME)-windows-x86_64.exe $(MAIN_FILE)

.PHONY: build-darwin
build-darwin:
	@CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 $(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(PROJECT_NAME)-darwin-aarch64 $(MAIN_FILE)
	@CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(PROJECT_NAME)-darwin-x86_64 $(MAIN_FILE)

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

.PHONY: test
test:
	go test ./...

.PHONY: archive
define create_archive
	cd $(BUILD_DIR) && tar -czvf $1-$2.tar.gz $(PROJECT_NAME)-$2$3 \
		-C ../ LICENSE  README.md \
		-C ./cmd/rules/ rules_example.json rules_example.yaml
	cd $(BUILD_DIR) && sha256sum $1-$2.tar.gz > $1-$2.tar.gz.sha256
endef
archive:
	$(call create_archive,$(PROJECT_NAME),linux-aarch64)
	$(call create_archive,$(PROJECT_NAME),linux-x86_64)
	$(call create_archive,$(PROJECT_NAME),windows-aarch64,.exe)
	$(call create_archive,$(PROJECT_NAME),windows-x86_64,.exe)
	$(call create_archive,$(PROJECT_NAME),darwin-aarch64)
	$(call create_archive,$(PROJECT_NAME),darwin-x86_64)

.PHONY: release
release: clean test build-all archive
