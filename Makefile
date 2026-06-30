BUILD_DIR ?= out
GO_BUILDFLAGS ?= -buildvcs=false -trimpath
LDFLAGS ?= -s -w -X main.version=$(shell cat VERSION 2>/dev/null || echo "dev")

.PHONY: all
all: build

$(BUILD_DIR):
	mkdir -p $(BUILD_DIR)

define build-binary
	CGO_ENABLED=0 GOOS=$(1) GOARCH=$(2) go build $(GO_BUILDFLAGS) -ldflags="$(LDFLAGS)" -o $(BUILD_DIR)/tapewrap-$(1)-$(2)$(3) ./cmd/tapewrap
	ln -sf tapewrap-$(1)-$(2)$(3) $(BUILD_DIR)/stow-$(1)-$(2)$(3)
endef

.PHONY: build
build: $(BUILD_DIR)
	CGO_ENABLED=0 go build $(GO_BUILDFLAGS) -ldflags="$(LDFLAGS)" -o $(BUILD_DIR)/tapewrap ./cmd/tapewrap
	ln -sf tapewrap $(BUILD_DIR)/stow

.PHONY: cross
cross: $(BUILD_DIR)
	$(call build-binary,linux,amd64,)
	$(call build-binary,linux,arm64,)
	$(call build-binary,darwin,amd64,)
	$(call build-binary,darwin,arm64,)
	$(call build-binary,windows,amd64,.exe)

.PHONY: clean
clean:
	rm -rf $(BUILD_DIR)

.PHONY: test
test:
	go test ./...
