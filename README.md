# tldr go client

## Usage

```
usage: tldr [OPTION]... SEARCH

available commands:
    -v, --version           print version and exit
    -h, --help              print this help and exit
    -u, --update            update local database
    -p, --platform PLATFORM select platform, supported are linux / osx / sunos / common
    -a, --list-all          list all available commands for the current platform
    -r, --render PATH       render a local page for testing purposes
```

![Example Output](https://raw.githubusercontent.com/mstruebing/tldr-go-client/master/docs/example.png "Example Output")

## Install

Just copy the executable anywhere on your system, preferably in some folder where 
your `$PATH` variable will find it.

If you want to build it yourself see below.

## Dependencies
At the current state you just need unzip installed on your system to unzip the pages, nothing else is needed.

To build it yourself you need golang installed.

## Building

If you want to build it yourself you can use the `Makefile` and type `make build`.
This will put the `tldr` bin in a `dist` folder.
If you want to compile it without it just do a `go build` in the root of this repository.

To install it on your system you can do a simple `make install` in the root of this repository.
This will build the executable file and copy it into `.local/bin`
Make sure you have this directory in your `$PATH`.
Otherwise you can build the executable yourself and copy it wherever you want. Or simply adjust the `Makefile` to your needs.


## Contribution

This is an early stage and maybe there are some bugs.
Contribution in form of issues, suggestions, testing and pull requests are very welcome.

If you contribute code wise please make sure to run `gofmt` and `govet`
