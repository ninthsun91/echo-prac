BINARY_NAME=myapp
.DEFAULT_GOAL := run

build:
	GOARCH=amd64 GOOS=linux go build -v -o ./target/$(BINARY_NAME)-linux ./cmd/myapp
	GOARCH=amd64 GOOS=darwin go build -v -o ./target/$(BINARY_NAME)-darwin ./cmd/myapp
	GOARCH=amd64 GOOS=windows go build -v -o ./target/$(BINARY_NAME).exe ./cmd/myapp

run: build
	$(eval UNAME_S := $(shell uname -s))
	$(eval UNAME_M := $(shell uname -m))
	@if [ "$(UNAME_S)" = "Linux" ]; then \
		./target/$(BINARY_NAME)-linux; \
	elif [ "$(UNAME_S)" = "Darwin" ]; then \
		./target/$(BINARY_NAME)-darwin; \
	elif [ "$(UNAME_S)" = "MINGW32_NT" -o "$(UNAME_S)" = "MINGW64_NT" -o "$(UNAME_S)" = "MSYS_NT" ]; then \
		./target/$(BINARY_NAME).exe; \
	else \
		echo "Unsupported OS"; \
	fi

test:
	@if [ -z "$(TEST)" ]; then \
		echo "Running all tests under ./internal/...\n"; \
		go test -v ./internal/...; \
	else \
		echo "Running tests under ./internal/... with pattern $(TEST)\n"; \
		go test -v ./internal/... -run $(TEST); \
	fi

test-e2e:
	@if [ -z "$(TEST)" ]; then \
		echo "Running all e2e tests...\n"; \
		go test -v ./tests/...; \
	else \
		echo "Running e2e tests with pattern $(TEST)\n"; \
		go test -v ./tests/... -run $(TEST); \
	fi

clean:
	go clean
	rm -rf ./target