.PHONY: build
build:
	go build -o build/ ./...

.PHONY: generate
generate:
	go run ./cmd/golox-ast/main.go

.PHONY: install
install:
	go install ./...

.PHONY: test
test:
	go test -v -race ./...
