// Package main provides ...
package main

import (
	"bufio"
	"flag"
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

const (
	listAllUsage  = "list all available commands for the current platform"
	platformUsage = "select platform; supported are: linux, osx, sunos, common"
	renderUsage   = "render a local page for testing purposes"
	updateUsage   = "update local database"
	versionUsage  = "print version and exit"
)

func printVersion() {
	fmt.Println("tldr v 1.0.2")
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
	return path.Join(homeDir, ".tldr")
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

func updateLocal() {
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
			if !strings.HasPrefix(lines[i+1], ">") {
				fmt.Println()
			}
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

func printPageInPath(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		log.Fatal("ERROR: Page doesn't exist")
	} else {
		lines, err := readLines(path)
		if err != nil {
			log.Fatal("ERROR: Something went wrong while reading the page")
		}
		printPage(lines)
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
		updateLocal()
	}

	version := flag.Bool("version", false, versionUsage)
	flag.BoolVar(version, "v", false, versionUsage)

	update := flag.Bool("update", false, updateUsage)
	flag.BoolVar(update, "u", false, updateUsage)

	render := flag.String("render", "", renderUsage)
	flag.StringVar(render, "r", "", renderUsage)

	listAll := flag.Bool("list-all", false, listAllUsage)
	flag.BoolVar(listAll, "a", false, listAllUsage)

	platform := flag.String("platform", "", platformUsage)
	flag.StringVar(platform, "p", "", platformUsage)

	flag.Parse()

	if *version {
		printVersion()
	} else if *update {
		updateLocal()
	} else if *render != "" {
		printPageInPath(*render)
	} else if *listAll {
		listAllPages()
	} else if *platform != "" {
		page := flag.Arg(0)
		if page == "" {
			log.Fatal("ERROR: no page provided")
		}
		printPageForPlatform(*platform, flag.Arg(0))
	} else {
		page := flag.Arg(0)
		if page == "" {
			flag.PrintDefaults()
			os.Exit(0)
		}
		printSinglePage(page)
	}
}
