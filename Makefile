BINARY_NAME := go-kitty

.PHONY: build
build:
	go build -o $(BINARY_NAME) .

clean:
	rm -f $(BINARY_NAME)