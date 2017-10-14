// Package main provides ...
package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path"
)

func fetchPages() {
	cacheDir := getCacheDir()
	cmd := exec.Command("git", "clone", "https://github.com/tldr-pages/tldr", path.Join(cacheDir, "tldr-git"))
	err := cmd.Run()
	if err != nil {
		log.Fatal("ERROR: Can't fetch tldr repository")
	}
}

func getHomeDirectory() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal("ERROR: " + err.Error())
	}
	if usr.HomeDir == "" {
		log.Fatal("ERROR: Can't load user's home folder path")
	}

	return usr.HomeDir
}

func getCacheDir() string {
	homeDir := getHomeDirectory()
	return path.Join(homeDir, ".cache", "tldr-go")
}

func createCacheDir() {
	cacheDir := getCacheDir()
	os.MkdirAll(cacheDir, 0755)
}

func copyPages() {
	cacheDir := getCacheDir()
	err := os.Rename(path.Join(cacheDir, "tldr-git", "pages"), path.Join(cacheDir, "pages"))
	if err != nil {
		log.Fatal("ERROR: " + err.Error())
	}
}

func removeCacheDir() {
	cacheDir := getCacheDir()
	os.RemoveAll(cacheDir)
}

func setup() {
	// TODO: read commit hash and put it into root of cache dir
	createCacheDir()
	fetchPages()
	copyPages()
}

func update() {
	// TODO: check for newer version via commit hash
	removeCacheDir()
	setup()
}

func main() {
	cacheDir := getCacheDir()
	if _, err := os.Stat(path.Join(cacheDir, "pages")); os.IsNotExist(err) {
		update()
	}

	fmt.Println("Hello world!")
}
