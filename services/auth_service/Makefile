# Run all tests
test:
	go test -v ./...

# Run tests with coverage
test-coverage:
	go test -v -cover ./...

# Run tests and ignore logs from Gin
test-silent:
	go test -v ./... | grep -v "\[GIN\]"
