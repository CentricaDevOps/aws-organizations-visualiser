NAME := "aws-organizations-visualiser"

build: test
	@echo "Building..."
	@go build -o bin/$(NAME) -v

run:
	@echo "Running..."
	@go run main.go

test: lint-check
	@echo "Testing..."
	@go test -v ./...

test-integration:
	@echo "Testing integration..."
	@INTEGRATION=true go test -v ./...

lint-check:
	@echo "Linting..."
	@go vet ./...

lint-fix:
	@echo "Linting and fixing..."
	@go fmt ./...