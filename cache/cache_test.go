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
	cacheDirectory, err := cacheDir()

	if err != nil || !strings.HasSuffix(cacheDirectory, ".cache/tldr") {
		t.Error("Expected to end with `.cache/tldr` but got", err, cacheDirectory)
	}

	os.Setenv("XDG_CACHE_HOME", "/tmp")
	cacheDirectory, err = cacheDir()

	if err != nil || (cacheDirectory != "/tmp/tldr") {
		t.Error("Expected to be `/tmp/.cache/tldr` but got", err, cacheDirectory)
	}

	os.Setenv("XDG_CACHE_HOME", "")
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
		t.Error("Expected to successfully pull a page from the cache")
	}

	_, err = r.Markdown("linux", "hostnamee")

	if err == nil {
		t.Error("Expected to return an error for non existing pages")
	}
}

func TestPages(t *testing.T) {
	remote := "https://tldr.sh/assets/tldr.zip"
	ttl := time.Hour * 24 * 7
	r, _ := NewRepository(remote, ttl)

	pages, err := r.Pages()

	if err != nil {
		t.Error("Expected to successfully retrieve all pages.")
	}

	if len(pages) == 0 {
		t.Error("Expected to find some pages")
	}
}

func TestHistory(t *testing.T) {

	repo := Repository{
		directory: "/tmp/.cache/tldr",
		remote:    "https://tldr.sh/assets/tldr.zip",
		ttl:       time.Hour * 24 * 7,
	}

	repo.makeCacheDir()

	if err2 := repo.RecordHistory("git-pull"); err2 != nil {
		t.Error("Expected to record history successful.")
	}

	repo.RecordHistory("git-pull")
	repo.RecordHistory("git-pull")
	repo.RecordHistory("git-push")
	repo.RecordHistory("git-push")
	repo.RecordHistory("git-fetch")
	repo.RecordHistory("git-pull")

	history, err := repo.loadHistory()
	if err != nil {
		t.Error("Expected to load history successful.")
	}

	length := len(*history)
	if length != 3 {
		t.Error("Expected to have 3 history records.")
	}

	rec1 := HistoryRecord{
		page:  "git-push",
		count: 2,
	}

	if (*history)[0] != rec1 {
		t.Errorf("Expected first record to be %+v", rec1)
	}

	rec2 := HistoryRecord{
		page:  "git-fetch",
		count: 1,
	}

	if (*history)[1] != rec2 {
		t.Errorf("Expected second record to be %+v", rec2)
	}

	rec3 := HistoryRecord{
		page:  "git-pull",
		count: 4,
	}

	if (*history)[2] != rec3 {
		t.Errorf("Expected third record to be %+v", rec3)
	}

}
