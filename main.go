// Package main provides ...
package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"path"
	"runtime"
	"strings"
)

func printHelp() {
	fmt.Println("usage: tldr [-v] [OPTION]... SEARCH")
	fmt.Println()
	fmt.Println("available commands:")
	fmt.Println("    -v, --version           print version and exit")
	fmt.Println("    -h, --help              print this help and exit")
	fmt.Println("    -u, --update            update local database")
	fmt.Println("    -p, --platform PLATFORM select platform, supported are linux / osx / sunos / common")
	fmt.Println("    -a, --list-all          list all available commands for the current platform")
	fmt.Println("    -r, --render PATH       render a local page for testing purposes")
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
	fmt.Println("fetching pages...")
	err := downloadFile(cacheDir+"/tldr.zip", "http://tldr-pages.github.io/assets/tldr.zip")
	if err != nil {
		log.Fatal(err.Error())
		log.Fatal("ERROR: Can't fetch tldr repository")
	}
}

func unzipPages() {
	cacheDir := getCacheDir()
	fmt.Println("unpacking pages...")
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
	fmt.Println("All done!")
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

func getSystems() []string {
	var systems []string
	pagesDir := getPagesDir()
	currentSystem := getCurrentSystem()
	systems = append(systems, currentSystem)
	systems = append(systems, "common")

	availableSystems, err := ioutil.ReadDir(path.Join(pagesDir))
	if err != nil {
		log.Fatal("ERROR: Something bad happened while reading directories")
	}

	for _, availableSystem := range availableSystems {
		if availableSystem.Name() != "index.json" && availableSystem.Name() != currentSystem && availableSystem.Name() != "common" {
			systems = append(systems, availableSystem.Name())
		}
	}

	return systems
}

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func listAllPages() {
	currentSystem := getCurrentSystem()
	pagesDir := getPagesDir()
	pages, err := ioutil.ReadDir(path.Join(pagesDir, currentSystem))
	if err != nil {
		log.Fatal("ERROR: Can't read pages for current platform: " + currentSystem)
	}

	for _, page := range pages {
		fmt.Println(page.Name()[:len(page.Name())-3])
	}
}

func convertExample(line string) string {
	var processedLine string = line
	const BLUE = "\x1b[34;1m"
	const RED = "\x1b[31;1m"
	processedLine = strings.Replace(processedLine, "{{", BLUE, -1)
	processedLine = strings.Replace(processedLine, "}}", RED, -1)
	return strings.Replace(processedLine, "`", "", -1)
}

func printPage(lines []string) {
	const GREEN = "\x1b[32;1m"
	const RED = "\x1b[31;1m"
	const RESET = "\x1b[30;1m"
	for i, line := range lines {
		if strings.HasPrefix(line, "#") {
			fmt.Println(line[2:])
			fmt.Println()
		}

		if strings.HasPrefix(line, ">") {
			fmt.Println(line[2:])
			fmt.Println()
		}

		if strings.HasPrefix(line, "-") {
			fmt.Printf("%s%s%s\n", GREEN, line, RESET)
			fmt.Printf("    %s%s%s\n", RED, convertExample(lines[i+2]), RESET)
			if i < len(lines)-3 {
				fmt.Println()
			}
		}
	}
}

func printPageForPlatform(platform string, page string) {
	pagesDir := getPagesDir()
	platformDir := path.Join(pagesDir, platform)
	file := platformDir + "/" + page + ".md"
	if _, err := os.Stat(file); os.IsNotExist(err) {
		log.Fatal("ERROR: no page found for " + page + " in platform " + platform)
	} else {
		lines, err := readLines(file)
		if err != nil {
			log.Fatal("ERROR: Something went wrong while reading the page")
		}
		printPage(lines)
	}
}

func printSinglePage(page string) {
	pagesDir := getPagesDir()
	systems := getSystems()

	for index, system := range systems {
		systemDir := path.Join(pagesDir, system)
		file := systemDir + "/" + page + ".md"
		if _, err := os.Stat(file); os.IsNotExist(err) {
			if index == 1 {
				log.Fatal("ERROR: no page found for " + page)
			}
		} else {
			lines, err := readLines(file)
			if err != nil {
				log.Fatal("ERROR: Something went wrong while reading the page")
			}
			printPage(lines)
			break
		}
	}
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
		printHelp()
	case "--help":
		printHelp()
	case "-v":
		printVersion()
	case "--version":
		printVersion()
	case "-u":
		update()
	case "--update":
		update()
	case "-a":
		listAllPages()
	case "--list-all":
		listAllPages()
	case "-p":
		if len(args) > 2 {
			printPageForPlatform(args[1], args[2])
		} else {
			log.Fatal("ERROR: No platform provided or page provided")
		}
	case "--platform":
		if len(args) > 2 {
			printPageForPlatform(args[1], args[2])
		} else {
			log.Fatal("ERROR: No platform provided or page provided")
		}
	default:
		printSinglePage(args[0])
	}
}
