.PHONY: build
build:
	go build -o build/ ./...

.PHONY: install
install:
	go install ./...

.PHONY: test
test:
	go test -v -race ./...
