BINARY_NAME=itodo
INSTALL_DIR=$(HOME)/.local/bin

.PHONY: build install clean demo publish-test publish

build:
	@go build -o $(BINARY_NAME) .
	@echo "Build complete: $(BINARY_NAME)"

install: build
	@mkdir -p $(INSTALL_DIR)
	@mv $(BINARY_NAME) $(INSTALL_DIR)/
	@chmod +x $(INSTALL_DIR)/$(BINARY_NAME)
	@echo "Installed: $(BINARY_NAME) -> $(INSTALL_DIR)/$(BINARY_NAME)"

clean:
	@rm -f $(BINARY_NAME)
	@echo "Cleaned: $(BINARY_NAME)"

demo: build
	@command -v vhs >/dev/null 2>&1 || (echo "Install VHS: go install github.com/charmbracelet/vhs@latest" && exit 1)
	@command -v ffmpeg >/dev/null 2>&1 || (echo "Install ffmpeg (e.g. brew install ffmpeg)" && exit 1)
	@vhs assets/demo.tape -o assets/demo.gif
	@echo "Created: demo.gif"

publish-test:
	@go mod tidy
	@goreleaser release --snapshot --clean
	@echo "Published test successful."

publish:
	@go mod tidy
	@goreleaser release --clean
	@echo "Published successful."
