// Package main provides ...
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/mstruebing/tldr"
	"github.com/mstruebing/tldr/cache"
)

// Help message constants
const (
	listAllUsage  = "list all available commands for the current platform"
	platformUsage = "select platform; supported are: linux, osx, sunos, common"
	renderUsage   = "render a local page for testing purposes"
	updateUsage   = "update local database"
	versionUsage  = "print version and exit"
)

const (
	remoteURL = "https://tldr.sh/assets/tldr.zip"
	ttl       = time.Hour * 24 * 7
)

func printVersion() {
	fmt.Println("tldr v 1.0.5")
	fmt.Println("Copyright (C) 2017 Max Strübing")
	fmt.Println("Source available at https://github.com")
}

func main() {
	repository, err := cache.NewRepository(remoteURL, ttl)
	if err != nil {
		log.Fatalf("ERROR: creating cache repository: %s", err)
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
		err = repository.Reload()
		if err != nil {
			log.Fatalf("ERROR: updating cache: %s", err)
		}
	} else if *render != "" {
		if _, err := os.Stat(*render); os.IsNotExist(err) {
			log.Fatal("ERROR: page doesn't exist")
		}

		page, err := os.Open(*render)
		if err != nil {
			log.Fatal("ERROR: opening the page")
		}
		defer page.Close()

		err = tldr.Write(page, os.Stdout)
		if err != nil {
			log.Fatalf("ERROR: rendering the page: %s", err)
		}
	} else if *listAll {
		pages, err := repository.Pages(tldr.CurrentPlatform())
		if err != nil {
			log.Fatalf("ERROR: getting pages: %s", err)
		}

		for _, page := range pages {
			fmt.Println(page)
		}
	} else if *platform != "" {
		page := flag.Arg(0)
		if page == "" {
			log.Fatal("ERROR: no page provided")
		}

		markdown, err := repository.Markdown(*platform, page)
		if err != nil {
			log.Fatalf("ERROR: getting markdown for '%s/%s': %s", *platform, page, err)
		}
		defer markdown.Close()

		err = tldr.Write(markdown, os.Stdout)
		if err != nil {
			log.Fatalf("ERROR: writing markdown: %s", err)
		}
	} else {
		page := flag.Arg(0)
		if page == "" {
			flag.PrintDefaults()
			os.Exit(0)
		}

		platform := tldr.CurrentPlatform()
		markdown, err := repository.Markdown(platform, page)
		if err != nil {
			var platforms []string
			platforms, err = tldr.AvailablePlatforms(repository)
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
	}
}
