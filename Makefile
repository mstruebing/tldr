COMPILE_COMMAND = go build -o bin/tldr cmd/tldr/main.go

# Test command
TEST = go test

# Set source dir and scan source dir for all go files
SRC_DIR = .
SOURCES = $(shell find $(SRC_DIR) -type f -name '*.go')

install: build 
	mkdir -p ~/.local/bin && cp bin/tldr ~/.local/bin

build: $(SOURCES)
	$(COMPILE_COMMAND)

.PHONY: clean
clean:
	rm -Rf bin && rm -Rf ~/.tldr
