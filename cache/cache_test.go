package cache

import (
	"os"
	"strings"
	"testing"
	"time"
)

func contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}

func TestCacheDir(t *testing.T) {
	cacheDir, err := cacheDir()

	if err != nil || !strings.HasSuffix(cacheDir, ".tldr") {
		t.Error("Expected linux, got ", err, cacheDir)
	}
}

func TestNewRepository(t *testing.T) {
	remote := "https://tldr.sh/assets/tldr.zip"
	ttl := time.Hour * 24 * 7
	r, _ := NewRepository(remote, ttl)
	cacheDir, _ := cacheDir()

	if r.directory != cacheDir {
		t.Error("Expected directory the same as cacheDir, got", r.directory)
	}

	if r.remote != remote {
		t.Error("Expected remote to be the same as called with, got", r.remote)
	}

	if r.ttl != ttl {
		t.Error("Expected ttl to be the same as called with, got", r.ttl)
	}

	_, err := os.Stat(cacheDir)

	if os.IsNotExist(err) {
		t.Error("Expected cache directory to extist but it isn't")
	}
}

func TestPlattforms(t *testing.T) {
	remote := "https://tldr.sh/assets/tldr.zip"
	ttl := time.Hour * 24 * 7
	r, _ := NewRepository(remote, ttl)

	platforms, _ := r.AvailablePlatforms()
	if len(platforms) != 5 {
		t.Error("Expected 5 Platforms, got", len(platforms))
	}

	if !contains(platforms, "linux") {
		t.Error("Expected linux in platforms, got", platforms)
	}

	if !contains(platforms, "common") {
		t.Error("Expected common in platforms, got", platforms)
	}

	if !contains(platforms, "osx") {
		t.Error("Expected osx in platforms, got", platforms)
	}

	if !contains(platforms, "sunos") {
		t.Error("Expected sunos in platforms, got", platforms)
	}

	if !contains(platforms, "windows") {
		t.Error("Expected windows in platforms, got", platforms)
	}
}

func TestReload(t *testing.T) {
	remote := "https://tldr.sh/assets/tldr.zip"
	ttl := time.Hour * 24 * 7
	r, _ := NewRepository(remote, ttl)

	err := r.Reload()

	if err != nil {
		t.Error("Expected to successfully reload the repository, got", err)
	}
}

func TestMarkdown(t *testing.T) {
	remote := "https://tldr.sh/assets/tldr.zip"
	ttl := time.Hour * 24 * 7
	r, _ := NewRepository(remote, ttl)

	_, err := r.Markdown("linux", "hostname")

	if err != nil {
		t.Error("Exptected to successfully pull a page from the cache")
	}

	_, err = r.Markdown("linux", "hostnamee")

	if err == nil {
		t.Error("Exptected to return an error for non existing pages")
	}
}

func TestPages(t *testing.T) {
	remote := "https://tldr.sh/assets/tldr.zip"
	ttl := time.Hour * 24 * 7
	r, _ := NewRepository(remote, ttl)

	pages, err := r.Pages()

	if err != nil {
		t.Error("Exptected to successfully retrieve all pages.")
	}

	if len(pages) == 0 {
		t.Error("Exptected to find some pages")
	}
}

func TestLoadFromRemote(t *testing.T) {
	remote := "https://tldr.sh/assets/tldr.zip"
	ttl := time.Hour * 24 * 7
	r, _ := NewRepository(remote, ttl)

	err := r.loadFromRemote()

	if err != nil {
		t.Error("Exptected to successfully retrieve all pages.")
	}
}
