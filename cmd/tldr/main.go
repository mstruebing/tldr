// Package main provides ...
package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"runtime"
	"time"

	"github.com/mstruebing/tldr"
	"github.com/mstruebing/tldr/cache"
)

// Help message constants
const (
	listAllUsage  = "list all available commands for the current platform"
	platformUsage = "select platform; supported are: linux, osx, sunos, common"
	pathUsage     = "render a local page for testing purposes"
	updateUsage   = "update local database"
	versionUsage  = "print version and exit"
	randomUsage   = "prints a random page"
	historyUsage  = "show the latest search history"
)

const (
	remoteURL = "https://tldr.sh/assets/tldr.zip"
	ttl       = time.Hour * 24 * 7
)

const currentPlattform = runtime.GOOS

func printVersion() {
	fmt.Println("tldr v 1.3.1")
	fmt.Println("Copyright (C) 2017 Max Strübing")
	fmt.Println("Source available at https://github.com/mstruebing/tldr")
}

func listAllPages() {
	repository, err := cache.NewRepository(remoteURL, ttl)
	if err != nil {
		log.Fatalf("ERROR: creating repository: %s", err)
	}

	pages, err := repository.Pages()
	if err != nil {
		log.Fatalf("ERROR: getting pages: %s", err)
	}

	for _, page := range pages {
		fmt.Println(page)
	}
}

func printPageInPath(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		log.Fatal("ERROR: page doesn't exist")
	}

	page, err := os.Open(path)
	if err != nil {
		log.Fatal("ERROR: opening the page")
	}
	defer page.Close()

	err = tldr.Write(page, os.Stdout)
	if err != nil {
		log.Fatalf("ERROR: rendering the page: %s", err)
	}
}

func printPage(page string) {
	if page == "" {
		flag.PrintDefaults()
		os.Exit(0)
	}

	repository, err := cache.NewRepository(remoteURL, ttl)
	if err != nil {
		log.Fatalf("ERROR: creating cache repository: %s", err)
	}

	platform := tldr.CurrentPlatform(currentPlattform)
	markdown, err := repository.Markdown(platform, page)
	if err != nil {
		var platforms []string
		platforms, err = tldr.AvailablePlatforms(repository, currentPlattform)
		if err != nil {
			log.Fatalf("ERROR: getting available platforms: %s", err)
		}

		for _, platform = range platforms {
			markdown, err = repository.Markdown(platform, page)
			if err == nil {
				break
			}
		}
		if err != nil {
			log.Fatalf("ERROR: no page found for '%s' in any available platform", page)
		}
	}
	defer markdown.Close()

	err = tldr.Write(markdown, os.Stdout)
	if err != nil {
		log.Fatalf("ERROR: writing markdown: %s", err)
	}

	err = repository.RecordHistory(page)
	if err != nil {
		log.Fatalf("ERROR: saving history: %s", err)
	}
}

func printPageForPlatform(page string, platform string) {
	if page == "" {
		log.Fatal("ERROR: no page provided")
	}

	repository, err := cache.NewRepository(remoteURL, ttl)
	if err != nil {
		log.Fatalf("ERROR: creating cache repository: %s", err)
	}

	markdown, err := repository.Markdown(platform, page)
	if err != nil {
		log.Fatalf("ERROR: getting markdown for '%s/%s': %s", platform, page, err)
	}
	defer markdown.Close()

	err = tldr.Write(markdown, os.Stdout)
	if err != nil {
		log.Fatalf("ERROR: writing markdown: %s", err)
	}
}

func printRandomPage() {
	repository, err := cache.NewRepository(remoteURL, ttl)
	if err != nil {
		log.Fatalf("ERROR: creating cache repository: %s", err)
	}

	pages, err := repository.Pages()
	if err != nil {
		log.Fatalf("ERROR: getting pages: %s", err)
	}
	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s) // initialize local pseudorandom generator
	printPage(pages[r.Intn(len(pages))])
}

func updatePages() {
	repository, err := cache.NewRepository(remoteURL, ttl)
	if err != nil {
		log.Fatalf("ERROR: creating cache repository: %s", err)
	}
	err = repository.Reload()
	if err != nil {
		log.Fatalf("ERROR: updating cache: %s", err)
	}
}

func printHistory() {
	repository, err := cache.NewRepository(remoteURL, ttl)
	if err != nil {
		log.Fatalf("ERROR: creating cache repository: %s", err)
	}

	history, err := repository.LoadHistory()
	if err != nil {
		log.Fatalf("ERROR: error loading history: %s", err)
	}

	hisLen := len(*history)
	if hisLen == 0 {
		fmt.Println("No history is available yet")
	} else { //default print last 10.
		size := int(math.Min(10, float64(hisLen)))
		for i := 1; i <= size; i++ {
			record := (*history)[hisLen-i]
			fmt.Printf("%s\n", record)
		}
	}
}

func main() {
	version := flag.Bool("version", false, versionUsage)
	flag.BoolVar(version, "v", false, versionUsage)

	update := flag.Bool("update", false, updateUsage)
	flag.BoolVar(update, "u", false, updateUsage)

	path := flag.String("path", "", pathUsage)
	// f like file
	flag.StringVar(path, "f", "", pathUsage)

	listAll := flag.Bool("list-all", false, listAllUsage)
	flag.BoolVar(listAll, "a", false, listAllUsage)

	platform := flag.String("platform", "", platformUsage)
	flag.StringVar(platform, "p", "", platformUsage)

	random := flag.Bool("random", false, randomUsage)
	flag.BoolVar(random, "r", false, randomUsage)

	history := flag.Bool("history", false, historyUsage)
	flag.BoolVar(history, "t", false, historyUsage)

	flag.Parse()

	if *version {
		printVersion()
	} else if *update {
		updatePages()
	} else if *path != "" {
		printPageInPath(*path)
	} else if *listAll {
		listAllPages()
	} else if *platform != "" {
		page := flag.Arg(0)
		printPageForPlatform(page, *platform)
	} else if *random {
		printRandomPage()
	} else if *history {
		printHistory()
	} else {
		page := flag.Arg(0)
		printPage(page)
	}
}
