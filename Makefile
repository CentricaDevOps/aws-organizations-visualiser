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
	@go fmt ./... | wc -l | grep 0 > /dev/null || (echo "Please run 'make lint-fix' to fix linting errors" && exit 1)

lint-fix:
	@echo "Linting and fixing..."
	@go fmt ./...