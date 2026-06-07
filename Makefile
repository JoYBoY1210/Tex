APP_NAME := tex
CMD_DIR  := cmd/tex
MAIN     := $(CMD_DIR)/main.go
BIN_DIR  := bin

.PHONY: run build clean

run:
	go run ./$(CMD_DIR)

build:
	mkdir -p $(BIN_DIR)
	go build -o $(BIN_DIR)/$(APP_NAME) ./$(CMD_DIR)

clean:
	rm -rf $(BIN_DIR)