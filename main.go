// Package main provides ...
package main

import (
	"fmt"
	"os/exec"
)

func fetchPages() {
	cmd := exec.Command("git", "clone", "https://github.com/tldr-pages/tldr", "/tmp/tldr-pages")
	err := cmd.Run()
	if err != nil {
		fmt.Println("ERROR: Can't fetch tldr repository")
	}
}

func main() {
	fmt.Println("Hello world!")
}
