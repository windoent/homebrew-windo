.PHONY: build clean build-all release sha256

BINARY_NAME=windo
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS=-s -w -X github.com/windoent/homebrew-windo/cmd.version=$(VERSION) -X github.com/windoent/homebrew-windo/cmd.buildTime=$(BUILD_TIME)

build:
	CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -o $(BINARY_NAME) .

build-all: clean
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -buildvcs=false -ldflags  "$(LDFLAGS)" -o $(BINARY_NAME)-darwin-arm64 .
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -buildvcs=false -ldflags "$(LDFLAGS)" -o $(BINARY_NAME)-darwin-amd64 .
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -buildvcs=false -ldflags "$(LDFLAGS)" -o $(BINARY_NAME)-windows.exe .
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -buildvcs=false -ldflags "$(LDFLAGS)" -o $(BINARY_NAME)-linux-amd64 .
	@echo "Build complete:"
	@ls -lh $(BINARY_NAME)-*

sha256:
	@echo "=== SHA256 Checksums ==="
	@echo "darwin-arm64: $$(sha256sum $(BINARY_NAME)-darwin-arm64 | awk '{print $$1}')"
	@echo "darwin-amd64: $$(sha256sum $(BINARY_NAME)-darwin-amd64 | awk '{print $$1}')"
	@echo "windows: $$(sha256sum $(BINARY_NAME)-windows.exe | awk '{print $$1}')"
	@echo "linux-amd64: $$(sha256sum $(BINARY_NAME)-linux-amd64 | awk '{print $$1}')"

release: build-all sha256
	@echo ""
	@echo "=== Next Steps ==="
	@echo "1. Create GitLab Release at: https://github.com/windoent/homebrew-windo/-/releases/new?tag=$(VERSION)"
	@echo "2. Upload the following files:"
	@ls -1 $(BINARY_NAME)-darwin-* $(BINARY_NAME)-windows.exe $(BINARY_NAME)-linux-* 2>/dev/null | sed 's/^/   - /'
	@echo "3. Update Formula/windo.rb with the SHA256 values above"
	@echo "4. Push tap update to homebrew-windo repository"

clean:
	rm -f $(BINARY_NAME)-*

install:
	CGO_ENABLED=0 go install -ldflags "$(LDFLAGS)" .

test:
	go test ./...

lint:
	golangci-lint run

fmt:
	go fmt ./...
