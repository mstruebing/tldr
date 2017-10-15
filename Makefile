COMPILE_COMMAND = go build -o dist/tldr main.go

# Test command
TEST = go test

# Set source dir and scan source dir for all go files
SRC_DIR = .
SOURCES = $(shell find $(SRC_DIR) -type f -name '*.go')

install: build 
	mkdir -p ~/.local/bin && cp dist/tldr ~/.local/bin

build: $(SOURCES)
	$(COMPILE_COMMAND)

.PHONY: clean
clean:
	rm -Rf dist && rm -Rf ~/.cache/tldr-go
