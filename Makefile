COMPILE_COMMAND = go build -o bin/tldr cmd/tldr/main.go

# Set source dir and scan source dir for all go files
SRC_DIR = .
SOURCES = $(shell find $(SRC_DIR) -type f -name '*.go')

BINARIES = $(wildcard bin/*)

build: $(SOURCES)
	$(COMPILE_COMMAND)

install: build
	mkdir -p ~/.local/bin && cp bin/tldr ~/.local/bin

build-all-binaries: $(SOURCES) clean
	# doesn't work on my machine and not in travis, see: https://github.com/golang/go/wiki/GoArm
	# GOOS=android GOARCH=arm $(COMPILE_COMMAND) && mv ./bin/tldr ./bin/tldr-android-arm
	# GOOS=darwin  GOARCH=arm $(COMPILE_COMMAND) && mv ./bin/tldr ./bin/tldr-darwin-arm
	# GOOS=darwin  GOARCH=arm64 $(COMPILE_COMMAND) && mv ./bin/tldr ./bin/tldr-darwin-arm64
	GOOS=darwin    GOARCH=386      $(COMPILE_COMMAND) && mv ./bin/tldr ./bin/tldr-darwin-386
	GOOS=darwin    GOARCH=amd64    $(COMPILE_COMMAND) && mv ./bin/tldr ./bin/tldr-darwin-amd64
	GOOS=dragonfly GOARCH=amd64    $(COMPILE_COMMAND) && mv ./bin/tldr ./bin/tldr-dragonfly-amd64
	GOOS=freebsd   GOARCH=386      $(COMPILE_COMMAND) && mv ./bin/tldr ./bin/tldr-freebsd-386
	GOOS=freebsd   GOARCH=amd64    $(COMPILE_COMMAND) && mv ./bin/tldr ./bin/tldr-freebsd-amd64
	GOOS=freebsd   GOARCH=arm      $(COMPILE_COMMAND) && mv ./bin/tldr ./bin/tldr-freebsd-arm
	GOOS=linux     GOARCH=386      $(COMPILE_COMMAND) && mv ./bin/tldr ./bin/tldr-linux-386
	GOOS=linux     GOARCH=amd64    $(COMPILE_COMMAND) && mv ./bin/tldr ./bin/tldr-linux-amd64
	GOOS=linux     GOARCH=arm      $(COMPILE_COMMAND) && mv ./bin/tldr ./bin/tldr-linux-arm
	GOOS=linux     GOARCH=arm64    $(COMPILE_COMMAND) && mv ./bin/tldr ./bin/tldr-linux-arm64
	GOOS=linux     GOARCH=ppc64    $(COMPILE_COMMAND) && mv ./bin/tldr ./bin/tldr-linux-ppc64
	GOOS=linux     GOARCH=ppc64le  $(COMPILE_COMMAND) && mv ./bin/tldr ./bin/tldr-linux-ppc64le
	GOOS=linux     GOARCH=mips     $(COMPILE_COMMAND) && mv ./bin/tldr ./bin/tldr-linux-mips
	GOOS=linux     GOARCH=mipsle   $(COMPILE_COMMAND) && mv ./bin/tldr ./bin/tldr-linux-mipsle
	GOOS=linux     GOARCH=mips64   $(COMPILE_COMMAND) && mv ./bin/tldr ./bin/tldr-linux-mips64
	GOOS=linux     GOARCH=mips64le $(COMPILE_COMMAND) && mv ./bin/tldr ./bin/tldr-linux-mips64le
	GOOS=netbsd    GOARCH=386      $(COMPILE_COMMAND) && mv ./bin/tldr ./bin/tldr-netbsd-386
	GOOS=netbsd    GOARCH=amd64    $(COMPILE_COMMAND) && mv ./bin/tldr ./bin/tldr-netbsd-amd64
	GOOS=netbsd    GOARCH=arm      $(COMPILE_COMMAND) && mv ./bin/tldr ./bin/tldr-netbsd-arm
	GOOS=openbsd   GOARCH=386      $(COMPILE_COMMAND) && mv ./bin/tldr ./bin/tldr-openbsd-386
	GOOS=openbsd   GOARCH=amd64    $(COMPILE_COMMAND) && mv ./bin/tldr ./bin/tldr-openbsd-amd64
	GOOS=openbsd   GOARCH=arm      $(COMPILE_COMMAND) && mv ./bin/tldr ./bin/tldr-openbsd-arm
	GOOS=plan9     GOARCH=386      $(COMPILE_COMMAND) && mv ./bin/tldr ./bin/tldr-plan9-386
	GOOS=plan9     GOARCH=amd64    $(COMPILE_COMMAND) && mv ./bin/tldr ./bin/tldr-plan9-amd64
	GOOS=solaris   GOARCH=amd64    $(COMPILE_COMMAND) && mv ./bin/tldr ./bin/tldr-solaris-amd64
	GOOS=windows   GOARCH=386      $(COMPILE_COMMAND) && mv ./bin/tldr ./bin/tldr-windows-386
	GOOS=windows   GOARCH=amd64    $(COMPILE_COMMAND) && mv ./bin/tldr ./bin/tldr-windows-amd64

compress-all-binaries: build-all-binaries
	for f in $(BINARIES); do      \
        tar czf $$f.tar.gz $$f;    \
    done
	@rm $(BINARIES)

test: $(SOURCES)
	@go test -v ./...
	@go tool vet .
	@test -z $(shell gofmt -s -l . | tee /dev/stderr) || (echo "[ERROR] Fix formatting issues with 'gofmt'" && exit 1)

.PHONY: clean
clean:
	rm -Rf bin && rm -Rf ~/.tldr
