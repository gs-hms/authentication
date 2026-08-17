.PHONY: fmt

fmt:
	gofmt -w .

vet:
	go vet ./...

lint:
	golangci-lint run

test:
	go test ./...

check:
	gofmt -l .
	go vet ./...
	golangci-lint run
	go test ./...

run:
	go run ./cmd/server/main.go
