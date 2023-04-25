.PHONY: build run test

build:
	go build -o manager

run:
	go run main.go

test:
	go test ./...

