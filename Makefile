COMPILE_COMMAND = go build -o dist/tldr main.go

# Test command
TEST = go test

# Set source dir and scan source dir for all go files
SRC_DIR = .
SOURCES = $(shell find $(SRC_DIR) -type f -name '*.go')

# Set test dir and scan test dir for all java files
TEST_DIR = src
TEST_SOURCES = $(shell find $(TEST_DIR) -type f -name '*_test.go')


# Targets
all: start

start: build
	./dist/tldr

build: $(SOURCES)
	$(COMPILE_COMMAND)

test: $(SOURCES) $(TEST_SOURCES)
	$(TEST)

.PHONY: clean
clean:
	rm -Rf dist

