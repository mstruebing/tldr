// Package main provides ...
package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"path"
	"runtime"
)

func printHelp() {
	fmt.Println("usage: tldr [-v] [OPTION]... SEARCH")
	fmt.Println()
	fmt.Println("available commands:")
	fmt.Println("    -v                   print verbose output")
	fmt.Println("    --version            print version and exit")
	fmt.Println("    -h, --help           print this help and exit")
	fmt.Println("    -u, --update         update local database")
	fmt.Println("    -c, --clear-cache    clear local database")
	fmt.Println("    -p, --platform=PLATFORM select platform, supported are linux / osx / sunos / common")
	fmt.Println("    -r, --render=PATH    render a local page for testing purposes")
}

func printVersion() {
	fmt.Println("tldr v 0.0.1")
	fmt.Println("Copyright (C) 2017 Max Strübing")
	fmt.Println("Source available at https://github.com")
}

func downloadFile(filepath string, url string) (err error) {
	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func fetchPages() {
	cacheDir := getCacheDir()
	err := downloadFile(cacheDir+"/tldr.zip", "http://tldr-pages.github.io/assets/tldr.zip")
	if err != nil {
		log.Fatal(err.Error())
		log.Fatal("ERROR: Can't fetch tldr repository")
	}
}

func unzipPages() {
	cacheDir := getCacheDir()
	cmd := exec.Command("unzip", cacheDir+"/tldr.zip", "-d", cacheDir)
	err := cmd.Run()
	if err != nil {
		log.Fatal("ERROR: Can't unzip pages")
	}

	os.Remove(cacheDir + "/tldr.zip")
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
	return path.Join(homeDir, ".tldr-go")
}

func getPagesDir() string {
	cacheDir := getCacheDir()
	return path.Join(cacheDir, "pages")
}

func createCacheDir() {
	cacheDir := getCacheDir()
	os.MkdirAll(cacheDir, 0755)
}

func removeCacheDir() {
	cacheDir := getCacheDir()
	os.RemoveAll(cacheDir)
}

func setup() {
	createCacheDir()
	fetchPages()
	unzipPages()
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

	if len(args) < 1 {
		printHelp()
		os.Exit(0)
	}

	switch args[0] {
	case "-h":
	case "--help":
		printHelp()
	case "-v":
	case "--version":
		printVersion()
	}

	currentSystem := getCurrentSystem()

	for index, folder := range []string{currentSystem, "common"} {
		systemDir := path.Join(pagesDir, folder)
		file := systemDir + "/" + args[0] + ".md"
		if _, err := os.Stat(file); os.IsNotExist(err) {
			if index == 1 {
				log.Fatal("ERROR: no page found for " + args[0])
			}
		} else {
			fmt.Println(file)
			break
		}
	}
}
