// Package main provides ...
package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"os/user"
)

func fetchPages() {
	cmd := exec.Command("git", "clone", "https://github.com/tldr-pages/tldr", "/tmp/tldr-pages")
	err := cmd.Run()
	if err != nil {
		fmt.Println("ERROR: Can't fetch tldr repository")
	}
}

func getHomeDirectory() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	if usr.HomeDir == "" {
		return "", errors.New("Can't load user's home folder path")
	}

	return usr.HomeDir, nil
}

func createCacheDir() error {
	homeDir, err := getHomeDirectory()
	if err != nil {
		os.Stderr.WriteString("ERROR: " + err.Error() + "\n")
		return err
	}

	os.MkdirAll(homeDir+"/.tldr-go", 755)
	return nil
}

func main() {
	createCacheDir()
	fmt.Println("Hello world!")
}
