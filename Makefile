BINARY_NAME := kitty

.PHONY: build
build:
	go build -o $(BINARY_NAME) .
