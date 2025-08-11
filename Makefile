BINARY=boilerplate-compose
BIN_DIR=bin

.PHONY: build
build:
	GO111MODULE=on go build -o $(BIN_DIR)/$(BINARY) ./cmd/boilerplate-compose

.PHONY: tidy
	go mod tidy