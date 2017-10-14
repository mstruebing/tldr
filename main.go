// Package main provides ...
package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path"
	"runtime"
)

func fetchPages() {
	gitDir := getGitDir()
	cmd := exec.Command("git", "clone", "https://github.com/tldr-pages/tldr", gitDir)
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

func getPagesDir() string {
	cacheDir := getCacheDir()
	return path.Join(cacheDir, "pages")
}

func getGitDir() string {
	cacheDir := getCacheDir()
	return path.Join(cacheDir, "tldr-git")
}

func createCacheDir() {
	cacheDir := getCacheDir()
	os.MkdirAll(cacheDir, 0755)
}

func copyPages() {
	gitDir := getGitDir()
	pagesDir := getPagesDir()
	err := os.Rename(path.Join(gitDir, "pages"), pagesDir)
	if err != nil {
		log.Fatal("ERROR: " + err.Error())
	}
}

func removeCacheDir() {
	cacheDir := getCacheDir()
	os.RemoveAll(cacheDir)
}

func setup() {
	createCacheDir()
	fetchPages()
	copyPages()
}

func update() {
	removeCacheDir()
	setup()
}

func getCurrentSystem() string {
	os := runtime.GOOS
	switch os {
	case "darwin":
		os = "osx"
	}

	return os
}

func getFoldersToSearch() []string {
	currentSystem := getCurrentSystem()
	return []string{currentSystem, "common"}
}

func main() {
	pagesDir := getPagesDir()
	if _, err := os.Stat(pagesDir); os.IsNotExist(err) {
		update()
	}

	args := os.Args[1:]
	currentSystem := getCurrentSystem()

	for index, folder := range []string{currentSystem, "common"} {
		systemDir := path.Join(pagesDir, folder)
		file := systemDir + "/" + args[0] + ".md"
		if _, err := os.Stat(file); os.IsNotExist(err) {
			if index == 2 {
				log.Fatal("ERROR: no page found for " + args[0])
			}
		}
	}

	fmt.Println("Hello world!")
}
