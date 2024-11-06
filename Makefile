.PHONY: build test clean run

# Build the application
build:
	go build -o bin/simulator cmd/simulator/main.go

# Run tests
test:
	go test ./...

# Clean build artifacts
clean:
	rm -rf bin/

# Run the application
run: build
	./bin/simulator

# Run tests with coverage
test-coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

# Format code
fmt:
	go fmt ./...

# Run linter
lint:
	golangci-lint run

# Generate mocks (if needed)
generate:
	go generate ./...