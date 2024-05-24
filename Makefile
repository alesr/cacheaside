.PHONY: tools
tools:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.58.1

.PHONY: lint
lint:
	@command -v golangci-lint > /dev/null || (echo "golangci-lint not found. Run 'make tools' to install it." && exit 1)
	@echo "Running lint..."
	@golangci-lint run
	@echo "Lint passed."

.PHONY: test
test:
	@echo "Running tests..."
	@go test -v -race ./...
	@echo "Tests passed."
