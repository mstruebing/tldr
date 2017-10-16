package tldr

import (
	"runtime"
)

// CurrentPlatform returns the platform name of the current system.
func CurrentPlatform() string {
	os := runtime.GOOS
	if os == "darwin" {
		os = "osx"
	}
	return os
}
