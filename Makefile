PORT ?= 8080

.PHONY: run test build clean

## Start the server
run:
	go run ./cmd/ledger -port=$(PORT)

## Run all tests
test:
	go test ./... -v -count=1 -race

## Build binary
build:
	go build -o bin/tiny-ledger ./cmd/ledger

## Remove build artifacts
clean:
	rm -rf bin/
