BINARY_NAME=adeleApp

build:
	@go mod vendor
	@echo "Building..."
	@go build -o tmp/${BINARY_NAME} .
	@echo "Build complete!"

run: build
	@echo "Starting..."
	@./tmp/${BINARY_NAME}

clean:
	@echo "Cleaning..."
	@go clean
	@rm -f tmp/${BINARY_NAME}
	@echo "Cleaned!"

start: run

stop:
	@echo "Stopping..."
	@-pkill -SIGTERM -f "./tmp/${BINARY_NAME}"
	@PID=$$(pgrep -f "./tmp/${BINARY_NAME}"); \
    while kill -0 $$PID 2>/dev/null; do sleep 0.1; done

restart: stop start

fmt:
	gofmt -s -w .

fmt-check:
	@test -z "$$(gofmt -s -l .)"

vet:
	go vet ./...

lint:
	golangci-lint run ./...

test:
	go test -race -cover ./...

# ci runs fmt-check, vet, and tests. Lint is intentionally omitted here —
# the GitHub Actions workflow invokes golangci-lint via the dedicated
# golangci-lint-action so this target stays usable without a local
# golangci-lint install. Run `make lint` directly if you want lint locally.
ci: fmt-check vet test
