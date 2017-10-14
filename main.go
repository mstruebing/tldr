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
	homeDir, err := getHomeDirectory()
	if err != nil {
		os.Stderr.WriteString("ERROR: " + err.Error() + "\n")
	}
	cmd := exec.Command("git", "clone", "https://github.com/tldr-pages/tldr", homeDir+"/.cache/tldr-go/pages-git")
	err = cmd.Run()
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

	os.MkdirAll(homeDir+"/.cache/tldr-go", 0755)
	return nil
}

func copyPages() {
	homeDir, err := getHomeDirectory()
	if err != nil {
		os.Stderr.WriteString("ERROR: " + err.Error() + "\n")
		return
	}

	err = os.Rename(homeDir+"/.cache/tldr-go/pages-git/pages", homeDir+"/.cache/tldr-go/pages")
	if err != nil {
		os.Stderr.WriteString("ERROR: " + err.Error() + "\n")
	}
}

func main() {
	createCacheDir()
	fetchPages()
	copyPages()
	err := createCacheDir()
	if err != nil {
		os.Exit(1)
	}
	fmt.Println("Hello world!")
}
