BINARY_NAME=itodo
INSTALL_DIR=$(HOME)/.local/bin

.PHONY: build install clean publish-test publish

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

publish-test:
	@go mod tidy
	@goreleaser release --snapshot --clean
	@echo "Published test successful."

publish:
	@go mod tidy
	@goreleaser release --clean
	@echo "Published successful."
