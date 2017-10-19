package tldr

import (
	"runtime"
)

// CommonPlatform is the common platform.
const CommonPlatform = "common"

// CurrentPlatform returns the platform name of the current system.
func CurrentPlatform() string {
	os := runtime.GOOS
	if os == "darwin" {
		os = "osx"
	}
	return os
}

// AvailablePlatforms returns all the available platforms that are supported.
func AvailablePlatforms(r Repository) ([]string, error) {
	platforms, err := r.AvailablePlatforms()
	if err != nil {
		return nil, err
	}

	current := CurrentPlatform()
	var currentFound, commonFound bool
	for _, p := range platforms {
		if p == current {
			currentFound = true
		} else if p == CommonPlatform {
			commonFound = true
		}
	}

	if !currentFound {
		platforms = append(platforms, CurrentPlatform())
	}
	if !commonFound {
		platforms = append(platforms, CommonPlatform)
	}
	return platforms, nil
}
