.PHONY: build run test lint clean

APP_NAME := ibookfs
BUILD_DIR := build

build:
	go build -o $(BUILD_DIR)/$(APP_NAME) ./cmd/server

run:
	go run ./cmd/server

test:
	go test -v -race -cover ./...

lint:
	golangci-lint run

clean:
	rm -rf $(BUILD_DIR)

tidy:
	go mod tidy

docker-build:
	docker build -t $(APP_NAME) .

docker-run:
	docker run -p 8080:8080 $(APP_NAME)
