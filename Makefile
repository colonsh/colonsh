BIN := colonsh
INSTALL_DIR := $(HOME)/bin

.PHONY: all build test install uninstall clean

all: build

## Build the binary
build:
	go build -o $(BIN) .

## Run all tests
test:
	go test ./...

## Build and install to $(HOME)/bin
## After installing, reload your shell: source ~/.zshrc
install: build
	mv ./$(BIN) $(INSTALL_DIR)/$(BIN)

## Remove the installed binary from $(HOME)/bin
uninstall:
	rm -f $(INSTALL_DIR)/$(BIN)

## Remove the built binary from the project directory
clean:
	rm -f $(BIN)
