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
	go test -v ./internal/...

test-e2e:
	go test -v ./tests/...

clean:
	go clean
	rm -rf ./target