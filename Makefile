build:
	@go build -o fs

lint:
	@golangci-lint run .
