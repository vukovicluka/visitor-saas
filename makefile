APP_NAME := visitor
BUILD_DIR := build

.PHONY: build build-linux clean run

build:
	go build -o $(BUILD_DIR)/$(APP_NAME) ./cmd/visitor

build-linux:
	GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(APP_NAME)-linux-amd64 ./cmd/visitor

clean:
	rm -rf $(BUILD_DIR)

run: build
	$(BUILD_DIR)/$(APP_NAME)
