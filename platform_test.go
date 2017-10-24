package tldr

import (
	"testing"
)

func TestCurrentPlattform(t *testing.T) {
	var currentPlattform string

	currentPlattform = CurrentPlatform("linux")
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
