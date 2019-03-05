package tldr

import (
	"github.com/mstruebing/tldr/cache"
	"testing"
	"time"
)

const (
	remoteURL = "https://tldr.sh/assets/tldr.zip"
	ttl       = time.Hour * 24 * 7
)

func testSliceEqual(a, b []string) bool {

	if a == nil && b == nil {
		return true
	}

	if a == nil || b == nil {
		return false
	}

	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

func TestCurrentPlattform(t *testing.T) {
	var currentPlattform string = CurrentPlatform("linux")

	if currentPlattform != "linux" {
		t.Error("Expected linux, got ", currentPlattform)
	}

	currentPlattform = CurrentPlatform("LINUX")
	if currentPlattform != "linux" {
		t.Error("Expected linux, got ", currentPlattform)
	}

	currentPlattform = CurrentPlatform("darwin")
	if currentPlattform != "osx" {
		t.Error("Expected osx, got ", currentPlattform)
	}

	currentPlattform = CurrentPlatform("DARWIN")
	if currentPlattform != "osx" {
		t.Error("Expected osx, got ", currentPlattform)
	}

	currentPlattform = CurrentPlatform("sunos")
	if currentPlattform != "sunos" {
		t.Error("Expected sunos, got ", currentPlattform)
	}

	currentPlattform = CurrentPlatform("SUNOS")
	if currentPlattform != "sunos" {
		t.Error("Expected sunos, got ", currentPlattform)
	}

	currentPlattform = CurrentPlatform("windows")
	if currentPlattform != "windows" {
		t.Error("Expected windows, got ", currentPlattform)
	}

	currentPlattform = CurrentPlatform("WINDOWS")
	if currentPlattform != "windows" {
		t.Error("Expected windows, got ", currentPlattform)
	}
}

func TestAvailablePlatforms(t *testing.T) {
	var availablePlatforms []string
	repository, _ := cache.NewRepository(remoteURL, ttl)

	availablePlatforms, _ = AvailablePlatforms(repository, "linux")
	if !testSliceEqual([]string{"common", "linux", "osx", "sunos", "windows"}, availablePlatforms) {
		t.Error("Expected to get all available platforms, got ", availablePlatforms)
	}

	availablePlatforms, _ = AvailablePlatforms(repository, "stuff")
	if !testSliceEqual([]string{"common", "linux", "osx", "sunos", "windows", "stuff"}, availablePlatforms) {
		t.Error("Expected to get all available platforms including 'stuff', got ", availablePlatforms)
	}
}
